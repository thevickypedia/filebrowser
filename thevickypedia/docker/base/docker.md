### Build
```shell
docker build -t filebrowser_base .
```

### Build (with upload)
```shell
export DOCKER_BUILDKIT=1
docker build \
  --secret id=git_token,src=<(printf "%s" "$GIT_TOKEN") \
  --secret id=repo_name,src=<(printf "%s" "$GITHUB_REPOSITORY") \
  -t filebrowser_base .
```

### Run
```shell
docker run -p 8080:80 filebrowser_base
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
  filebrowser_base
```

> Refer [this][stackoverflow] for docker port spec

### Copy Executable
```shell
docker cp "$(docker ps -aqf 'ancestor=filebrowser_base')":/filebrowser .
```

### Post-cleanup
```shell
docker stop $(docker ps -aqf "ancestor=filebrowser_base")
docker rm $(docker ps -aqf "ancestor=filebrowser_base")
docker rmi $(docker images -q "filebrowser_base") -f
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
