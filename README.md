# ws-moderator

## How to build

* install npm, go, podman/docker, podman-compose/docker-compose
* `make gui`
* `make moderator`

## How to configure

* update `gui/.env` if you want to allow remote access to the `moderator` from `gui` javascript

## How to run

* mongodb

```sh
cd docker
podman-compose up
```

* web app

```sh
make run
```

* your app will run at `http://localhost`
* mongo express runs at `http://localhost:8081`
