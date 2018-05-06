# Reddit Watcher

Watches specific subreddits for items matching a specific pattern.

## Usage

### Running the Bot

Running the bot is easy!  You just need a Telegram token that you can get from the [BotFather](https://core.telegram.org/bots#3-how-do-i-create-a-bot).

Once you have that, start up your server with the following command:

```bash
docker run -v `pwd`/config:/config stjohnjohnson/reddit-watcher:latest --token ${TELEGRAM_TOKEN}
```

In this example, I'm running the container with settings being saved to a local directory.

### Using the Bot

The bot monitors `/r/mechmarket` for specific patterns.  You can specify what you are looking for by the following commands:

#### `/buying <keyword>`

This will tell the bot to look for items matching that keyword that are being sold.  Sold means the listing includes "cash" or "paypal" in the "want" field.

#### `/help`

Replies with a simple help message listing all the available commands.

#### `/stop <keyword>`

This tells the bot to no longer look for that keyword.

#### `/items`

Outputs a list of your keywords and the number of matches found so far.

#### `/stats`

Outputs some simple stats about the bot.

## Todo

Here are some of the features I still want to add to this bot:

 - [ ] Vendor searches - Looking for specific vendors posts
 - [ ] Artisan searches - Looking for specific artisan posts
 - [ ] Group Buy alerts - Looking for any group buys
 - [ ] Giveaway alerts - Looking for anything matching "giveaway"
 - [ ] Multi-Region - Customize region beyond US
 - [ ] RegExp support - Be able to provide Regex instead of just strings
