### Pre-cleanup
```shell
rm -rf filebrowser filebrowser.db frontend/node_modules
rm -rf frontend/dist && mkdir -p frontend/dist && touch frontend/dist/.gitkeep
```

### Build
> :bulb: This step can be replicated for all Operating Systems by modifying the `DockerfileScratch`
```shell
docker build --no-cache --progress=plain -f DockerfileScratch -t filebrowser .
```

### Copy Executable
```shell
docker cp $(docker ps -aqf "ancestor=filebrowser"):/opt/filebrowser/filebrowser .
```

### Run
```shell
docker run -p 8080:8080 filebrowser
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
