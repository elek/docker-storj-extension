
This repository contains a docker desktop extension. In the current form this is nothing more just a UI for starting a Docker Registry <-> Storj gateway.

The started container includes a modified version of [distribution/distribution](https://github.com/distribution/distribution). It's a Docker registry which stores all the data and metadata on [Storj Decentralized cloud](https://storj.io) with implementing a new Storj driver.

This extension just starts the required gateway (and explain how it can be used), but later it can be improved to use full OAuth workflow to get STORJ access grant instead of require it in the form.

The source code of the gateway can be found [here](https://github.com/elek/distribution/tree/storj-driver) and can be tested without this extension with:

```
docker run -p 9999:5000 -e REGISTRY_STORAGE_STORJ_BUCKET=docker -e REGISTRY_STORAGE_STORJ_ACCESSGRANT=$(cat /tmp/grant) ghcr.io/elek/distribution:618d19fb
```

## Try it out

As the extension is not in the marketplace, you need:

 * Docker Desktop: see the [official documentation](https://docs.docker.com/desktop/install/linux-install/) to install it,  
 * Docker Extension CLI: Download it from [here](https://github.com/docker/extensions-sdk/releases/tag/v0.2.4) and copy it to `/usr/lib/docker/cli-plugins/` (or similar folder)

Extension can be installed with the command:

```
docker extension install ghcr.io/elek/docker-storj-extension:20220728-2
```

At first time, you need to create a new bucket and generate the access grant (you can use Storj web interface or uplink). 

For example on EU1 gateway you can do it at [https://eu1.storj.io/access-grants](https://eu1.storj.io/access-grants)

After starting the container, you can use the local registry backed by Storj:

Pull an image from Docker Hub:

```
docker pull elek/herbstag

Using default tag: latest
latest: Pulling from elek/herbstag
a0d0a0d46f8b: Pull complete
Digest: sha256:0d86302fb863deb9a1e7371a947e7d3c029c1d55a28bb5f20f7e09317254e562
Status: Downloaded newer image for elek/herbstag:latest
```

Tag the image with the new location:

```
docker tag elek/herbstag localhost:9999/elek/herbstag
```

And push it:

```
docker push localhost:9999/elek/herbstag

Using default tag: latest
The push refers to repository [localhost:9999/elek/herbstag]
e2eb06d8af82: Layer already exists
latest: digest: sha256:0d86302fb863deb9a1e7371a947e7d3c029c1d55a28bb5f20f7e09317254e562 size: 528
```

You can double check the content of the Storj bucket:

```
export UPLINK_ACCESS=$(cat /tmp/grant)    

uplink ls -r sj://docker/  | grep herbstag
OBJ     2022-07-28 11:36:51    71         docker/registry/v2/repositories/elek/herbstag/_manifests/revisions/sha256/0d86302fb863deb9a1e7371a947e7d3c029c1d55a28bb5f20f7e09317254e562/link
OBJ     2022-07-28 11:36:51    71         docker/registry/v2/repositories/elek/herbstag/_manifests/tags/latest/index/sha256/0d86302fb863deb9a1e7371a947e7d3c029c1d55a28bb5f20f7e09317254e562/link
OBJ     2022-07-28 11:36:52    71         docker/registry/v2/repositories/elek/herbstag/_manifests/tags/latest/current/link
OBJ     2022-06-30 11:17:15    71         docker/registry/v2/repositories/elek/herbstag/_layers/sha256/a0d0a0d46f8b52473982a3c466318f479767577551a53ffc9074c9fa7035982e/link
OBJ     2022-06-30 11:17:21    71         docker/registry/v2/repositories/elek/herbstag/_layers/sha256/833c7a986ed965eec8fe864223920c366fb0a25dd23edd0bdd2a4428fd0ce1e2/link
```

## Development

New images can be built and pushed in one step:

```
docker buildx build --push -t ghcr.io/elek/docker-storj-extension:20220728-1 . 
```

For local test run it's enough to use `--load`:

```
docker buildx build --load -t ghcr.io/elek/docker-storj-extension:20220728-1 .
```

Extension can be updated with:

```
docker extension update ghcr.io/elek/docker-storj-extension:20220728-3
```

Dev console can be turned on by:

```
docker extension dev debug ghcr.io/elek/docker-storj-extension
```