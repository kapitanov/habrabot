# Habrabot

A telegram bot that publishes items from an RSS feed into a Telegram channel.

Written for [Habr.com](https://habr.com/) but can be used for any RSS feed

## Build and run

1. Obtain a Telegram bot token (see [here](https://core.telegram.org/bots/api#authorizing-your-bot))
2. Create a Telegram channel to post messages into.
3. Write a `.env` file:

   ```env
   TELEGRAM_TOKEN=__place_telegram_token_here__
   TELEGRAM_CHANNEL=__place_telegram_channel_name_or_id_here__
   RSS_FEED=__place_rss_feed_url_here__
   ```

   See `example.env` for more details.

4. Run the following command:

   ```bash
   docker-compose up -d --build
   ```
