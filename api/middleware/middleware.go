// Package middleware defines all HTTP middleware functions
// used to before handlers of routes of the pmd-dx-api.
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/janek64/pmd-dx-api/api/db"
	"github.com/janek64/pmd-dx-api/api/handler"
	"github.com/janek64/pmd-dx-api/api/logger"
	"github.com/julienschmidt/httprouter"
)

// ResourceListParams checks for possible arguments of resource list queries, parses their
// values and stores them in a struct which is added to the context of the request.
func ResourceListParams(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Retrieve the parameters from the request
		queryParams := r.URL.Query()
		// Generate the ResourceListParams struct and add it to the context
		var params handler.ResourceListParams
		// sorting
		sort := queryParams.Get("sort")
		// Check if the value is one of the sort types
		if sort == db.IDAsc || sort == db.IDDesc || sort == db.NameAsc || sort == db.NameDesc {
			params.Sort.SortEnabled = true
			params.Sort.SortType = db.SortType(sort)
		} else {
			// Invalid ordering types are ignored instead of being answered with an error
			params.Sort.SortEnabled = false
		}
		// pagination
		var err error
		// If page is zero, set to default value
		if params.Pagination.PerPage, err = strconv.Atoi(queryParams.Get("per_page")); err != nil || params.Pagination.PerPage == 0 {
			params.Pagination.PerPage = 50
		}
		// If per_page is zero, set to default value
		if params.Pagination.Page, err = strconv.Atoi(queryParams.Get("page")); err != nil || params.Pagination.Page == 0 {
			params.Pagination.Page = 1
		}
		ctx := context.WithValue(r.Context(), handler.ResourceListParamsKey, params)
		// Call the handler with the created context
		h(w, r.WithContext(ctx), ps)
	}
}

// FieldLimitingParams checks for the "fields" argument of the query used for field limiting,
// parses the value and stores it in a struct which is added to the context of the request.
func FieldLimitingParams(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Retrieve the parameters from the request
		fields := strings.Split(r.URL.Query().Get("fields"), ",")
		// Generate the FieldLimitingParams struct and add it to the context
		var fieldLimitParams handler.FieldLimitingParams
		// Check if at least one value was provided
		if len(fields) > 0 && fields[0] != "" {
			fieldLimitParams.FieldLimitingEnabled = true
			fieldLimitParams.Fields = fields
		} else {
			fieldLimitParams.FieldLimitingEnabled = false
		}
		ctx := context.WithValue(r.Context(), handler.FieldLimitingParamsKey, fieldLimitParams)
		// Call the handler with the created context
		h(w, r.WithContext(ctx), ps)
	}
}

// LogRequest logs the request with the logger package by using a custom http.ResponseWriter.
func LogRequest(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		responseRecorder := logger.ResponseRecorder{ResponseWriter: w, Status: 200, Size: 0}
		h(&responseRecorder, r, ps)
		err := logger.LogRequest(r, responseRecorder)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Writing to the access log failed: %v", err)
		}
	}
}
