import {createDockerDesktopClient} from '@docker/extension-api-client';

export function Service() {

  const ddClient = createDockerDesktopClient();

  async function runDockerInfo() {
    const result = await ddClient.docker.cli.exec('info', [
      '--format',
      '"{{json .}}"',
    ]);
    return result.parseJsonObject();
  }

  return this;
}
