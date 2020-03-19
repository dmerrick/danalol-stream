## Tripbot loves you :robot: :heart:

[![GoDoc](https://godoc.org/github.com/dmerrick/tripbot?status.svg)](https://godoc.org/github.com/dmerrick/tripbot)
[![Go Report Card](https://goreportcard.com/badge/github.com/dmerrick/tripbot)](https://goreportcard.com/report/github.com/dmerrick/tripbot)
[![Build Status](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fdmerrick%2Ftripbot%2Fbadge&style=flat)](https://actions-badge.atrox.dev/dmerrick/tripbot/goto)
[![License: CC BY-NC-SA 4.0](https://img.shields.io/badge/License-CC%20BY--NC--SA%204.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-sa/4.0/)

This is the source code to [whereisdana.today](http://whereisdana.today), a 24/7 interactive [slow-tv](https://en.wikipedia.org/wiki/Slow_television) art project.

If you like it, please consider subscribing to my channel on [Twitch.tv](https://www.twitch.tv/ADanaLife_).
Thanks for watching!

-Dana ([dana.lol](https://dana.lol))


### How it all works

There are two main components.
There's the chatbot itself, which listens for user commands, and a VLC-based video server, which manages the currently-playing video.
They can be run on separate computers.

The general flow of information looks like this:

![](assets/infra-diagram.png)

*Not pictured: a relational database, an MPD-based audio server, and a NAS.*


### Running tripbot locally

You can use `docker-compose` to run tripbot on your own machine.
It is configured to spin up all of the dependencies for the project.
A helper script ([`bin/devenv`](https://github.com/dmerrick/tripbot/blob/master/bin/devenv)) has been created to make the process a little easier.
For example:

```bash
# (optional) create alias for devenv script
alias devenv="$(pwd)/bin/devenv"

# spin up tripbot stack on current machine
devenv up --daemon
# see running containers
devenv ps

# see logs for a specific container
devenv logs tripbot

# shut down everything
devenv down
```


### Other Useful Docs

#### Infra

See [infra/README.md](infra/README.md) for infra setup instructions.

#### Database

See [db/README.md](db/README.md) for database instructions.


