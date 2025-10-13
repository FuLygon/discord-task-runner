# Discord Tasker Runner

[![Publish to GitHub Container Registry](https://github.com/FuLygon/discord-task-runner/actions/workflows/publish-package.yaml/badge.svg)](https://github.com/FuLygon/discord-task-runner/actions/workflows/publish-package.yaml)

Configurable Discord Bot that can execute [Tasker](https://tasker.joaoapps.com) tasks remotely via FCM.

You'll need to go through Tasker's [Remote Action Execution](https://tasker.joaoapps.com/userguide/en/fcm.html) guide on how to set up your device to receive FCM messages for task execution, and retrieve the device token before setting up the Bot.

## Setting up Bot

### Docker Installation

- Prepare and set up the `config.yaml` file:

```bash
wget https://raw.githubusercontent.com/FuLygon/discord-task-runner/refs/heads/main/config.example.yaml -O config.yaml
```

- You'll also need a **Google Cloud Service Account File**. Mount both the config and service account file, then deploy service with docker. Compose file example:

```yaml
services:
  discord-tasker-runner:
    image: ghcr.io/fulygon/discord-task-runner:latest
    volumes:
      - ./config.yaml:/app/config.yaml
      - ./service-account.json:/app/service-account.json
```

### Source Installation

Ensure [Go](https://go.dev/doc/install) is installed.

- Clone the repo:

```bash
git clone https://github.com/FuLygon/discord-task-runner.git
cd discord-task-runner
```

- Prepare and set up the `config.yaml` file:

```bash
cp config.example.yaml config.yaml
```

- You'll also need a **Google Cloud Service Account File**. Rename service account JSON file to `service-account.json` and copy it to the same location as the `config.yaml` file. Then build and run:

```bash
go build -o discord-tasker-runner ./cmd/main.go
./discord-tasker-runner
```

## Known issues
- If the slash command returns `Unknown Integration` error, restarting your **Discord client** will most likely fix it. I haven't found any workaround for this issue.
- Tasker relies on **FCM** to execute task remotely, therefore there is no direct way to tell if the device successfully received the message and executed the task. If the device is not connected to the internet, the message will be **queued** by FCM until the device is back online. To mitigate this, you can configure the `ttl` in the config file for the slash command to set the lifespan of the message in the queue.