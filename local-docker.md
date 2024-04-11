### Cleanup
```shell
rm -rf filebrowser filebrowser.db frontend/node_modules
rm -rf frontend/dist && mkdir -p frontend/dist && touch frontend/dist/.gitkeep
```

### Build
```shell
docker build --no-cache --progress=plain -f LinuxDockerfile -t filebrowser .
```

[//]: # (### Run)

[//]: # (```shell)

[//]: # (docker run -p 8080:8080 filebrowser)

[//]: # (```)
