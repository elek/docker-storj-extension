import {createDockerDesktopClient} from '@docker/extension-api-client';

export class Service {

    constructor() {
        this.ddClient = createDockerDesktopClient();
    }

    static async StartContainer(bucket, grant) {
        const result = await this.ddClient.docker.cli.exec('run', [
            '--d',
            '-e',
            'REGISTRY_STORAGE_STORJ_ACCESSGRANT=' + grant,
            '-e',
            'REGISTRY_STORAGE_STORJ_BUCKET=' + bucket,
            '-p',
            '9999:5000',
            'ghcr.io/elek/distribution:618d19fb',
        ]);
        console.log(result.parseJsonObject());
        alert("Registry has been started")
    }
}

export default Service
