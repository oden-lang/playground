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

First build and push a new tag to the Docker Hub.

```bash
docker build -t owickstrom/oden-playground:<VERSION> .
docker push owickstrom/oden-playground:<VERSION>
```

Then SSH into the AWS server and run `~/update-oden-playground.sh`.

```bash
bash update-oden-playground.sh
Which tag of the Oden playground would you like to run?
<VERSION>
...
```
