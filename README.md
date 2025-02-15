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

### Developing locally

Set the following environment variables:

```bash
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="host=localhost port=5432 user=gowebdev password=gowebdevsecret dbname=gowebdevdb sslmode=disable"
export GOOSE_MIGRATION_DIR=migrations
export JWT_SECRET_KEY=jwtsecret
```

To run the code locally, first run an instance of the local database using `docker-compose`. If you don't have `docker` or `docker-compose`, its pretty easy to install them so go ahead and do that - you can use [this link](https://docs.docker.com/desktop/). Then, once you have both installed (which you can check by running `docker version` and `docker compose version`), you can run the following command to start a local postgres database which will have its data persisted:

```bash
docker compose up -d
```

The `-d` flag is to run it in detached mode - without it, all the logs will appear in your terminal and you will have start a new terminal session to run further commands. It's useful to learn about `docker` and `docker compose` so you understand how to build images and manage containers locally. You can leave this postgres database running, but if you ever want to stop it, you can run `docker compose down`.

## Dev Log

I didn't clone my backend template because I wasn't super happy with the implementation, especially the database migration files which were unnecessary. So, I slowly copied parts over, which required installing the following package:

```bash
go get github.com/jackc/pgx/v5
```

## License

This project is licensed under the MIT License.
Feel free to customize the content further as needed!
