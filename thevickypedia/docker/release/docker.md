### Build
```shell
docker build -t filebrowser_release .
```

### Build (alpine)
```shell
docker build -f Dockerfile.alpine --build-arg VERSION=v2.48.20 -t filebrowser_alpine .
```

### Run
```shell
docker run -p 8080:80 filebrowser_release
```

### Compose
```shell
docker-compose -f docker-compose.yml up
```

### Run with volume attached
```shell
docker run -d \
  --name filebrowser \
  -p 8080:80 \
  -v ${DOCKER_VOLUME_STORAGE:-$HOME}:/data \
  -v ${DOCKER_VOLUME_STORAGE:-$HOME}/.filebrowser/config:/config \
  --restart unless-stopped \
  filebrowser
```

> Refer [this][stackoverflow] for docker port spec

### Copy Executable
```shell
docker cp "$(docker ps -aqf 'ancestor=filebrowser')":/filebrowser .
```

### Post-cleanup
```shell
docker stop $(docker ps -aqf "ancestor=filebrowser")
docker rm $(docker ps -aqf "ancestor=filebrowser")
docker rmi $(docker images -q "filebrowser") -f
docker image prune -a
docker builder prune
```

<details>
<summary><strong>Overall Docker Cleanup</strong></summary>

> :warning: Deletes all containers, images and build cache

```shell
docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)
docker rmi $(docker images -q) -f
docker image prune -a
docker builder prune
```
</details>

[stackoverflow]: https://stackoverflow.com/a/62125889
