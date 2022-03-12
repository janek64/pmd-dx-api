// Package cache contains all definitions for interactions with
// the redis instance of the pmd-dx-api, consisting of type
// definitions and functions for connecting to redis and
// executing commands.
package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
)

// RedisConnectionError - type for redis connection error.
type RedisConnectionError struct {
	MissingVar string
}

// Error - implementation of the error interface.
func (e *RedisConnectionError) Error() string {
	return fmt.Sprintf("connecting to redis failed because of missing environment variable '%v'", e.MissingVar)
}

// CacheMissError- type for cache miss errors.
type CacheMissError struct {
	MissingKey string
}

// Error - implementation of the error interface.
func (e *CacheMissError) Error() string {
	return fmt.Sprintf("no redis cache entry for key '%v'", e.MissingKey)
}

// redisClient is the global client connection to the redis instance.
var redisClient *redis.Client

// InitRedis connects to the redis instance and sets the global redisClient variable.
func InitRedis() error {
	// Get connection data from environment
	redisURL, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		return &RedisConnectionError{"REDIS_URL"}
	}
	redisPassword, ok := os.LookupEnv("REDIS_PASSWORD")
	if !ok {
		return &RedisConnectionError{"REDIS_PASSWORD"}
	}
	// Connect to redis instance
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: redisPassword,
		DB:       0,
	})
	// Perform test ping
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	// Set eviction policy to delete least frequently used keys
	_, err = redisClient.ConfigSet(context.Background(), "maxmemory-policy", "allkeys-lfu").Result()
	if err != nil {
		return err
	}
	return nil
}

// CloseRedis closes the connection to the redis instance.
func CloseRedis() error {
	if redisClient == nil {
		return errors.New("no redis connection to close")
	}
	err := redisClient.Close()
	if err != nil {
		return err
	}
	redisClient = nil
	return nil
}

// responseHash represents a response entry in the redis cache
// and is used for scanning redis results.
type responseHash struct {
	HeaderBytes []byte `redis:"header"`
	Json        []byte `redis:"json"`
}

// GetCachedResponse fetches the redis cache entry for the url as the key
// and returns the decoded http.Header and json. If no entry is found, a
// CacheMissError will be returned.
func GetCachedResponse(url string) (http.Header, []byte, error) {
	if redisClient == nil {
		return nil, nil, errors.New("redis connection not initialized")
	}
	// Read the hash from redis: HMGET <url> header json
	readResult := redisClient.HMGet(context.Background(), url, "header", "json")
	// Store the data into an intermediate struct
	var result responseHash
	if err := readResult.Scan(&result); err != nil {
		return nil, nil, err
	}
	// If both byte slices are empty, a cache miss occurred
	if len(result.HeaderBytes) == 0 && len(result.Json) == 0 {
		return nil, nil, &CacheMissError{url}
	}
	// Deserialize []byte header to http.Header
	var header http.Header
	buffer := bytes.NewBuffer(result.HeaderBytes)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&header)
	if err != nil {
		return nil, nil, err
	}
	return header, result.Json, nil
}

// CacheResponseRecorder is a custom http.ResponseWriter recording the
// json/body and the status code of a HTTP response for caching purposes.
type CacheResponseRecorder struct {
	http.ResponseWriter
	Json   []byte
	Status int
}

// Write - implementation of http.ResponseWriter interface storing the body/json.
func (c *CacheResponseRecorder) Write(b []byte) (int, error) {
	c.Json = b
	return c.ResponseWriter.Write(b)
}

// WriteHeader - implementation of http.ResponseWriter interface storing the status code.
func (c *CacheResponseRecorder) WriteHeader(status int) {
	c.Status = status
	c.ResponseWriter.WriteHeader(status)
}

// StoreResponse stores the header and json of a HTTP response in the redis
// cache, using the URL as the key.
func StoreResponse(url string, header http.Header, json []byte) error {
	if redisClient == nil {
		return errors.New("redis connection not initialized")
	}
	// Serialize the http.Header to []byte to store it
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(header)
	if err != nil {
		return err
	}
	// Store the values as Hash in redis: HSET <url> header <header> json <json>
	redisClient.HSet(context.Background(), url, "header", buffer.Bytes(), "json", json)
	return nil
}
