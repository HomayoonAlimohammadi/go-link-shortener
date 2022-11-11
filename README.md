# A Link Shortener built with go

This app digests a long unusable and forgettable url into a
tiny token which can be used in order to redirect to the
original url.

## How to use

First of all you need to have `.env` and `.cobra.yaml` files for the
configuration and environment varialbles. The `.env` file is only used in
`docker-compose`, so if you're not gonna use it, ignore the `.env` file.

- Both files should be located in the root of the project (where this `README.md` exists).

Here is a template for both the `.cobra.yaml` and `.env` file:

- .cobra.yaml

```yaml
author: "Homayoon"
license: "Apache"
postgres:
  user: "postgres"
  password: "postgres"
  host: "postgres"
  port: 5432
  name: "linkshortener"
  sslmode: "disable"
```

- .env

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=linkshortener
DATABASE_HOST=postgres
```

If you want to use this without `Dockerfile` or `docker-compose`, don't forget to change the `postgres.host` to `localhost` or your own postgres remote host.

Run the application by:

```bash
make start
```

This should take care of downloading the dependencies and starting the project

Or if you decided to run the project in a containerized environment, run:

```bash
docker compose up --build
```

Which should run the application on `port 8000` alongside a postgres container.

## Enjoy!
