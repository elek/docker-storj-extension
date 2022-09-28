
This project is looking for new contributors. Please use standard Github workflow (opening issues, creating PRs).  

Contributions to this project are released under [Apache][LICENSE].

## Development

Extension has two parts: frontend and backend. They can be started locally, even if this running mode is slightly different from the real (Docker Desktop Extension) running environment. 


Before starting the frontend, please change the following line in `web/src/api.js`:

From:

```
const Api = new DockerDesktop()
```

To:

```
const Api = new DevServer()
```

Next, you can start the frontend: 

```
cd web
#only first time
npm install

npx vite
```


You can open the frontend at `http://localhost:3000`.

During local development the backend calls are forwarded to `http://127.0.0.1:5555/` (see `src/web/wite.config.js`). 

Backend can be started with `go`:

```
go run ./main.go run --config.api.address localhost:5555 
```

You are good if you see the big `Start registry button`.

## Testing the extension

Running as an extension os almost the same. You shouldn't change the `const Api` (or you should change it back) and install extension after creating a new image.

```
docker buildx build --load -t ghcr.io/elek/docker-storj-extension:1.0.0-nightly-0 .
docker extension install ghcr.io/elek/docker-storj-extension:1.0.0-nightly-0
```

After the first install, it's enough to upgrade the extension. Which can be done with the following very-simple helper script. (`36` is a unique number. You can start with `1` and increment, just to make sure you get the latest version to be deployed).

```
./dev-upgrade.sh 36
```

Dev console can be turned on by:

```
docker extension dev debug ghcr.io/elek/docker-storj-extension
```
