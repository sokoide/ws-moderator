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

## Note

* At the end of conversation, the history is saved as `$email.gob` in the current directory
* When you resume the conversation as the same email, it'll be used
* If it becomes big or fails for some reason, please delete the `$email.gob` and refresh the client browser to start a new conversation w/o a history
