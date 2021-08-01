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

Consider the following scenario:

- you have three repositories on two different platforms (github, gitlab)
- you would like to get notification on pull-requests on all of them into a slack / discord channel

You would have to do the following on all three of them:

- create a secret environment variable which holds the necessary token for all three of them
- create a webhook or a github action for the github repository
- create a pipeline action for the gitlab repository

Now, image you change the channel name... You'll have to update it for all three repositories. Imagine you have 20.

Enter Krok. You have to register all repositories with Krok. Affiliate the Discord Sender command with all three repositories
and set up the channel name in a shared setting including the token to discord and set that up as a command setting.

Now, once you have to change the channel name, you only have to change it once which will then be used by all three.

# Development

Developing Krok is fairly easy as the database that it requires is bootstrapped and ready made for you.
The integration tests create their own contained databases on each run, so they can be executed as many times as needed
and they won't step on each other's toe.

To create a test database run:

```
make test-db
```

To tear it down:

```
make rm-test-db
```

To build Krok, simply run:

```
make
```

This will build all binaries.

Fork it, work, work, zug, zug, create PR, done. The PR checker will run all tests but it's always advised to run them locally as well with

```
make test
```

Always have an accompanying issue with your PR so we know what problem it's trying to solve. Whether that is simply a refactor, or test coverage
increase, or even a typo fix, all PRs are welcomed and appreciated.

# Contributions

## Frontend

We would love to have help with Krok especially on the frontend side. None of us ( currently ) know frontend development too well.
We do plan on having a frontend for Krok alongside the CLI([krokctl](https://github.com/krok-o/krokctl)).

## Commands

Commands are never enough. Krok lives for the command it can execute on each webhook action. Without the commands Krok is almost useless.
Writing commands extends the universe of Krok and makes Krok itself more useful.

## Krok core

The server, which is this repository, can always use implementation of another platform for example. Or just general improvement of the code.