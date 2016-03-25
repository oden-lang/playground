# Oden Playground

## Run

You need [gb](https://getgb.io/) to build the server.

```bash
gb build
ODEN_CLI=path/to/bin/oden bin/playground
# or with docker-compose:
docker-compose build web
docker-compose up web
```

## Release

Requires access to the Heroku account.

```bash
heroku docker:release
```
