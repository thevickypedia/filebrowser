### Pre-cleanup
```shell
rm -rf filebrowser filebrowser.db frontend/node_modules
rm -rf frontend/dist && mkdir -p frontend/dist && touch frontend/dist/.gitkeep
```

### Build
```shell
docker build --no-cache --progress=plain -f DockerfileScratch -t filebrowser .
```

### Run
```shell
docker run -p 8080:80 filebrowser
```

> Refer [this][stackoverflow] for docker port spec

### Copy Executable
```shell
docker cp $(docker ps -aqf "ancestor=filebrowser"):/opt/filebrowser/filebrowser .
```

### Post-cleanup
```shell
docker stop $(docker ps -aqf "ancestor=filebrowser")
docker rm $(docker ps -aqf "ancestor=filebrowser")
docker rmi $(docker images -q "filebrowser") -f
docker builder prune
```

<details>
<summary><strong>Dangerous Step</strong></summary>

> :warning: Delete all containers, images and build cache

```shell
docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)
docker rmi $(docker images -q) -f
docker builder prune
```
</details>

[stackoverflow]: https://stackoverflow.com/a/62125889
