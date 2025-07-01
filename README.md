## Getting started

Before running the application you will need a working PostgreSQL installation and a valid DSN (data source name) for
connecting to the database.

Or simply run `docker-compose up` to start a PostgreSQL instance with a default database and user.

If you want to use temporal, check out the [**installation guide**](https://docs.temporal.io/cli#install) and start the temporal dev server with `temporal server start-dev `

Make sure that you're in the root of the project directory, fetch the dependencies with `go mod tidy`, then run the
application using `go run ./cmd/server`:

```
$ go mod tidy
$ go run ./cmd/server
```

If you make a request to the `GET /status` endpoint using `curl` you should get a response like this:

```
$ curl -i localhost:4444/status
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 09 May 2022 20:46:37 GMT
Content-Length: 23

{
    "Status": "OK",
}
```

## Project structure

Everything in the codebase is designed to be editable. Feel free to change and adapt it to meet your needs.

|                        |                                                            |
|------------------------|------------------------------------------------------------|
| **`assets`**           | Contains the non-code assets for the application.          |
| `↳ assets/migrations/` | Contains SQL migrations.                                   |
| `↳ assets/efs.go`      | Declares an embedded filesystem containing all the assets. |

|                    |                                                              |
|--------------------|--------------------------------------------------------------|
| **`cmd/server`**   | Api server specific code                                     |
| **`cmd/genmodel`** | Tool to generate database models (code) from database schema |

|                        |                                                                                 |
|------------------------|---------------------------------------------------------------------------------|
| **`internal`**         | Contains various helper packages used by the application.                       |
| `↳ internal/database/` | Contains your database-related code (setup, connection, queries, transactions). |
| `↳ internal/core/`     | Contains domain specific application logic.                                     |
| `↳ internal/funcs/`    | Contains helper functions to work with collections.                             |
| `↳ internal/version/`  | Contains the application version number definition.                             |
| `↳ internal/config/`   | Contains application configuration.                                             |

## Configuration settings

Configuration settings are managed via environment variables, with the environment variables read into your application
in the `run()` function in the `main.go` file.

You can try this out by setting a `HTTP_PORT` environment variable to configure the network port that the server is
listening on:

```
$ export APP__HTTP_PORT="9999"
$ go run ./cmd/server
```

## Managing SQL migrations and database model generation

The `Makefile` in the project root contains commands to easily create and work with database migrations:

|                                                |                                                                                                            |
|------------------------------------------------|------------------------------------------------------------------------------------------------------------|
| `$ make migrations/new name=add_example_table` | Create a new database migration in the `assets/migrations` directory.                                      |
| `$ make migrations/up`                         | Apply all up migrations.                                                                                   |
| `$ make migrations/down`                       | Apply all down migrations.                                                                                 |
| `$ make migrations/goto version=N`             | Migrate up or down to a specific migration (where N is the migration version number).                      |
| `$ make migrations/force version=N`            | Force the database to be specific version without running any migrations.                                  |
| `$ make migrations/version`                    | Display the currently in-use migration version.                                                            |
| `$ make model/gen`                             | Generate database models from database schema, generated file will be put in `internal/database/dbmodels`. |

These `Makefile` tasks are simply wrappers around calls to the `github.com/golang-migrate/migrate/v4/cmd/migrate` tool.
For more information, please see
the [official documentation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).

By default all 'up' migrations are automatically run on application startup using embeded files from
the `assets/migrations` directory. You can disable this by setting the `APP__DB__AUTO_MIGRATE` environment variable
to `false`.

## Admin tasks

The `Makefile` in the project root contains commands to easily run common admin tasks:

|                     |                                                                                           |
|---------------------|-------------------------------------------------------------------------------------------|
| `$ make tidy`       | Format all code and imports and tidy the `go.mod` file.                                   |
| `$ make audit`      | Run `go vet`, `staticheck`, `govulncheck`, execute all tests and verify required modules. |
| `$ make test`       | Run all tests.                                                                            |
| `$ make test/cover` | Run all tests and outputs a coverage report in HTML format.                               |
| `$ make build`      | Build a binary for the `cmd/api` application and store it in the `/tmp/bin` directory.    |
| `$ make run`        | Build and then run a binary for the `cmd/api` application.                                |
| `$ make docs`       | Generate api documentations in `pkg/docs` directory.                                      |

## Application version

The application version number is generated automatically based on your latest version control system revision number.
If you are using Git, this will be your latest Git commit hash. It can be retrieved by calling the `version.Get()`
function from the `internal/version` package.

Important: The version control system revision number will only be available when the application is built
using `go build`. If you run the application using `go run` then `version.Get()` will return the string `"unavailable"`.
