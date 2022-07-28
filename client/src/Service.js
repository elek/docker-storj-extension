import {createDockerDesktopClient} from '@docker/extension-api-client';

class Desktop {

    constructor() {
        this.ddClient = createDockerDesktopClient();
    }

    async StartContainer() {
        const result = await this.ddClient.docker.cli.exec('start', [
            'storj-gateway'
        ]);
        console.log(result)
    }

    async StopContainer() {
        const result = await this.ddClient.docker.cli.exec('stop', [
            'storj-gateway'
        ]);
        console.log(result)
    }

    // can be 'missing', 'stopped', 'running'
    async CheckStatus() {
        return await this.ddClient.docker.cli.exec('inspect', [
            'storj-gateway'
        ]).then(function (res) {
            const parsed = JSON.parse(res.stdout);
            console.log(parsed)
            if (parsed[0].State.Running) {
                return "running"
            }
            return "stopped"
        }, function (err) {
            return "missing"
        })
    }

    async CreateContainer(bucket, grant) {
        const result = await this.ddClient.docker.cli.exec('create', [
            '--name',
            'storj-gateway',
            '-e',
            'REGISTRY_STORAGE_STORJ_ACCESSGRANT=' + grant,
            '-e',
            'REGISTRY_STORAGE_STORJ_BUCKET=' + bucket,
            '-p',
            '9999:5000',
            'ghcr.io/elek/distribution:618d19fb',
        ]);
    }
}

class Mock {
    constructor() {
        this.status = "missing"
    }

    async StartContainer() {
        console.log("starting container")
        this.status = "running"
    }

    async StopContainer() {
        this.status = "stopped"

    }

    async CheckStatus() {
        return this.status
    }

    async CreateContainer(bucket, grant) {
        console.log("container is created")

    }
}

const Service = new Desktop()

export {Service as Service}
