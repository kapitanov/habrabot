# Habrabot

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/kapitanov/habrabot/Go)
![Latest release](https://img.shields.io/github/release/kapitanov/habrabot)
![License](https://img.shields.io/github/license/kapitanov/habrabot)

A telegram bot that publishes items from an RSS feed into a Telegram channel.

Written for [Habr.com](https://habr.com/) but can be used for any RSS feed.

## How to run

First, you would need to register a bot for Telegram using [`@BotFather`](https://t.me/BotFather).
See [here](https://core.telegram.org/bots/api#authorizing-your-bot) for more details.

Then, create a Telegram channel to post messages into.
Note that you may use the same bot for multiple channels.

The, you should add your bot to the channel as an administrator.

Finally, create a `.env` file:

```shell
TELEGRAM_TOKEN=my-telegram-bot-token
TELEGRAM_CHANNEL=@MyAwesomeChannel
RSS_FEED=https://habr.com/ru/rss/all/
```

Here:

* `TELEGRAM_TOKEN` is the token you obtained from `@BotFather`, e.g. `1234567890:AABBCCdde-ffGGHHiiJJkkLLmmNNooPPqqRR`.
* `TELEGRAM_CHANNEL` is the name of the channel you created, e.g. `@MyAwesomeChannel`.
* `RSS_FEED` is an URL of the RSS feed you want to publish, e.g. `https://habr.com/ru/rss/all/`.

See `example.env` for more details.

Now, you may create a `docker-compose.yaml` file that will instruct Docker how to run the bot:

```yaml
version: '3'
services:
    habrabot:
        image: ghcr.io/kapitanov/habrabot:latest
        restart: always
        env_file: ./.env
        logging:
            driver: "json-file"
            options:
                max-size: "10m"
                max-file: "1"
        volumes:
          - ./data/:/data
```

At this point, you should have the following files in your working directory:

```shell
$ ls -1a
docker-compose.yaml
.env
```

Now, you may run the bot using the following command:

```shell
docker compose up -d
```
## License

[MIT](LICENSE)
