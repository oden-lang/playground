# Oden Playground

## Run

```bash
ODENC=path/to/bin/odenc go run server.go
# or with docker-compose:
docker-compose build web
docker-compose up web
```

## Release

Requires access to the Heroku account.

```bash
heroku docker:release
```
