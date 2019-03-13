# CTF Flagbot

<img height="150" align="right" src="images/flagbot.png">

This is a bot created for managing CTF Discord guilds by [Helithumper](https://github.com/helithumper). The bot automatically removes CTF flags. In CTF competitions, the ability to block all flags from being sent in messages is crucial to running a fair competition. This bot will remove messages containing flag patterns and send a witty GIF in response.

## Usage

```plaintext
> flagbot -t <bot token> -c <configuration path>
<bot token> is the Discord bot's api token as given by the [developer console](https://discordapp.com/developers).
<configuration path> is the path to the configuration directory. The configuration directory contains the following

configuration
|> responses.txt is a plaintext file containing possible responses the bot will say upon a flag match
|> gifs.txt      is a plaintext file containing links to GIFs the bot will post after removing a flag
|> patterns.txt  is a plaintext file containing Regex patterns to match flags
```
