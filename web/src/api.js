import axios from "axios";

class Mock {

    constructor() {
        this.status = "MISSING"
        this.bucket = ""
    }

    getConfiguration() {
        return {
            "bucket": this.bucket
        }
    }

    configure(bucket, grant) {
        this.status = "STOPPED"
        return Promise.resolve("asd")
    }

    getStatus() {
        return this.status
    }

    create() {
        this.status = "STOPPED"
    }

    start() {
        this.status = "START"
    }

    localImages() {
        return [
            {
                image: "elek/herbstag",
            }
        ]
    }

    remoteImages() {
        return [
            {
                image: "elek/herbstag",
            }
        ]
    }

    push(image) {

    }

    pull(image) {

    }
}

class DevServer {
    configure(bucket, grant) {
        return axios.post("/api/v0/configure", {bucket: bucket, grant: grant})
    }

    status() {
        return axios.get("/api/v0/status")
    }

    start() {
        return axios.post("/api/v0/start")
    }

    stop() {
        return axios.post("/api/v0/stop")
    }

    localImages() {
        return axios.get("/api/v0/images/local")
    }

    push(image, tag) {
        return axios.post("/api/v0/push", {image: image, tag: tag})
    }
}

class DockerDesktop {
    configure(bucket, grant) {
        return axios.post("/api/v0/configure", {bucket: bucket, grant: grant})
    }

    status() {
        return window.ddClient.extension.vm.service.get("/api/v0/test")
    }

    start() {
        return window.ddClient.extension.vm.service.post("/api/v0/start")
    }

    stop() {
        return window.ddClient.extension.vm.service.post("/api/v0/stop")
    }

    localImages() {
        return window.ddClient.extension.vm.service.get("/api/v0/images/local")
    }

    push(image, tag) {
        return axios.post("/api/v0/push", {image: image, tag: tag})
    }
}

const Api = new DevServer()

export {Api as Api}