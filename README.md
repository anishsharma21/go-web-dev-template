# go-web-dev-template

Welcome to this repo. The purpose of this repo is to cover the necessary boilerplate for new Go web development projects. It builds on top of the [go-backend-starter-template](https://github.com/anishsharma21/go-backend-starter-template) and extends it with HTMX and other web related technologies/features. It includes the following technologies:

- **PostgreSQL Database**: Integrated with a PostgreSQL database.
- **Docker**: Containerised with a Dockerfile.
- **Goose**: Uses Goose for database migration handling.
- **Docker Compose**: Uses Docker Compose for local development setup.
- **Air**: Supports hot module reloading with Air.
- **HTMX**: For reactivity on the client side

## Getting Started

To get started with this project, clone the repository and follow the instructions below (more instructions will be added soon).

```bash
git clone https://github.com/anishsharma21/go-web-dev-template.git
cd go-web-dev-template
```

Then, run the following command to install all the go dependencies:

```bash
go mod download
```

## Developing locally

Begin by setting the following environment variables:

```bash
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="host=localhost port=5432 user=gowebdev password=gowebdevsecret dbname=gowebdevdb sslmode=disable"
export GOOSE_MIGRATION_DIR=migrations
export JWT_SECRET_KEY=jwtsecret
```

### Docker + postgres

To run the code locally you'll need to spin up a local `postgres` database instance. If you don't have `docker` install it using [this link](https://docs.docker.com/desktop/). Then, you should have `docker` installed (which you can check by running `docker version` and `docker compose version`), you can run the following command to start the database with persistent data which remain even after you close it:

```bash
docker compose up -d
```

The `-d` flag is to run it in detached mode - without it, all the logs will appear in your terminal and you will have start a new terminal session to run further commands. It's useful to learn about `docker` and `docker compose` so you understand how to build images and manage containers locally. You can leave this postgres database running, but if you ever want to stop it, you can run `docker compose down`.

### Hot Module Reloading (`air`)

This tool isn't required but its a quality of life / dx booster. `air` is used for Hot Module Reloading (HMR), which enables your code to automatically recompiled and re-run when changes are made:

```bash
go install github.com/air-verse/air@latest
```

The configuration for `air` is already present in the `.air.toml` file so you can simply run the command `air` on its own from the root of the project, and your server will be started up with HMR:

```bash
air
```

### Local Database migrations (`goose`)

Use the following command to install `goose` locally as it will not be included in the project as a dependency:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

With your database running in the background from the previous `docker compose` instructions, check that `goose` is correctly connected to your database by running the following command:

```bash
goose status
```

Ensure your database is running, then, run the following command to run the migration up:

```bash
goose up
```

If the migration went well, you should see `OK` messages next to each applied sql file, and the final line should say `successfully migrated database to version: ...`. You can check the status again to confirm the migrations occurred successfully. Further migration files can be created using the following command:

```bash
goose create {name of migration} sql
```

With the database running, run the following command to run the migration down:

```bash
goose down
```

### Other Development Instruction

When working within templates or handlers that render them, make sure to update the `selectors.go` file. This file contains CSS selectors for each template reduce the amount of hard coded values and duplication in the code.

Tests run locally use the local postgres database. To replicate the CICD environment, you can clear your database before running the tests. Use the following command to run tests locally:

```bash
go test ./tests -v
```

## License

This project is licensed under the MIT License.
Feel free to customize the content further as needed!
