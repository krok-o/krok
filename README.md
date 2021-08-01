# Krok

[![Coverage Status](https://coveralls.io/repos/github/krok-o/krok/badge.svg)](https://coveralls.io/github/krok-o/krok)
[![Discord Chat](https://img.shields.io/discord/799290432567771146.svg)](https://discord.gg/)

The main server of the hook management system.

# What is Krok?

Krok is a multi-platform webhook handling system which runs specified commands on certain, configured actions.

This is quite the mouthful so let's break it down.

### Multi-platform webhook handling system

Krok can handle webhooks for various implement platforms. At the time of this writing the following platforms are supported:

- [x] Github
- [x] Gitlab
- [ ] Gitea
- [ ] BitBucket
- [ ] ...

Once a repository is created (more on that in the [How do I use it?](#how-do-i-use-it) section) it will receive webhook for the
configured events from these platforms, i.e.: push, pull, pull-request, issue comment, etc. Whatever the respective platform supports.

### Specific commands

There is a whole section on what [commands](#commands) are. To tl;dr they are small, runnable scripts which perform something on
the respective action. Like, sending a slack message, running some command, performing code generation, building a hugo blog... etc.

# How do I use it?

To set up and use krok, refer to the Krok documentation [Installation](https://krok.app/basics/installation/) section. There, you'll
find Configuration as well.

Once Krok is up-and-running, refer to the [Tutorials](https://krok.app/tutorial/) section for some usage scenarios.

TL;DR:

Krok can run as a server somewhere and listen for hook events like a Bot for Github webhooks. In order to do that, the following
has to happen:

Register a thing called a `Repository`. This Repository contains information about the webhook. Where it is, what the callback
url is, and what kind of events it's subscribed too. Those events are platform specific.

Then, take this repository and affiliate commands to it. By adding relationships to commands you specify what commands
should execute on the event that happens. Currently, we don't support running specific commands per event, but maybe in the
future we'll do that. In any case, once the action is set up, and an event happens, Krok runs these commands and passes over
certain details to the command so it can perform the action it's supposed to. The command can do whatever an executing binary
is capable of, which is virtually limitless as long as the necessary credentials are provided.

Krok can save these securely and pass them along to the command so it can perform its function.

# Scenarios / Use cases

# Development

# Contributions

## Frontend

We would love to have help with Krok especially on the frontend side. None of us ( currently ) know frontend development too well.
We do plan on having a frontend for Krok alongside the CLI([krokctl](https://github.com/krok-o/krokctl)).

## Commands

Commands are never enough. Krok lives for the command it can execute on each webhook action. Without the commands Krok is almost useless.
Writing commands extends the universe of Krok and makes Krok itself more useful.