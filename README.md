# Oden Playground

## Run

```bash
ODEN_CLI=path/to/bin/oden go run server.go
# or with docker-compose:
docker-compose build web
docker-compose up web
```

## Release

Requires access to the Heroku account.

```bash
heroku docker:release
```
