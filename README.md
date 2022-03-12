# pmd-dx-api

This project offers a RESTful API with data from **Pokémon Mystery Dungeon: Rescue Team DX** in JSON format. It was inspired by the [PokeAPI](https://github.com/PokeAPI/pokeapi) project and is half for my fun, half for practising Go and other technologies.

Technologies and core libraries used in this project:
* [Go](https://go.dev/)
* [net/http](https://pkg.go.dev/net/http)
* [httprouter](https://github.com/julienschmidt/httprouter)
* [PostgreSQL](https://www.postgresql.org/) with [pgx](https://github.com/jackc/pgx)
* [redis](https://redis.io/) with [go-redis/redis](https://github.com/go-redis/redis/)
* [lumberjack](https://github.com/natefinch/lumberjack)

The data is provided as .csv-files ready to be imported into the database with the provided scripts. It was gratefully (manually) collected from [serebii.net](https://serebii.net), [bulbapedia.bulbagarden.net](https://bulbapedia.bulbagarden.net) and [game8.co](https://game8.co).

Pokémon and Pokémon character names are trademarks of Nintendo.