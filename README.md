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

Start up Krok, but note that the `hookbase` value has to be set to your public IP address so the hook creation can work. Otherwise it will fail, because
localhost is not supported as a callback url.

```
➜  krok git:(main) ✗ ./bin/darwin/amd64/krok --file-vault-location .tmp --hostname 0.0.0.0:9998 --plugin-location .tmp/plugins --hookbase <your-ip>:9998
8:56AM INF Please set a global secret key... Randomly generating one for now...
8:56AM INF Start listening...

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.1.17
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:9998
```

Once everything is running, use [krokctl](https://github.com/krok-o/krokctl) to test your running server.

You should have a test repository on some platform so Krok can create a webhook for it. This guide uses Github.

First, you'll need to save a token:

```
➜  krokctl git:(main) ✗ ./bin/darwin/amd64/krokctl create vcs --token <token> --vcs 1
Success!
```

Second, you'll need a command to test with. A couple test commands can be found under the [plugins](https://github.com/krok-o/plugins) repository.
Download some, build them, and upload the tar-ed up binary after making sure you exported the necessary credentials.

```
➜  krokctl git:(main) ✗ export KROK_API_KEY_ID=api-key-id
➜  krokctl git:(main) ✗ export KROK_API_KEY_SECRET=secret
➜  krokctl git:(main) ✗ export KROK_EMAIL=admin@admin.com
./bin/darwin/amd64/krokctl upload command --file slack-notification.tar.gz
ID      NAME                    HASH                                                                    LOCATION        FILENAME              SCHEDULE        ENABLED REPOSITORIES
1       slack-notification      5bac4e2aeff81e2453dd099128eaa65cd076d02bad7b6927e5afd4892cbacc2a        .tmp/plugins    slack-notification                    true
```

If the upload succeeded, you should see the above table.

Finally, register a repository:

```
./bin/darwin/amd64/krokctl create repository --events ping --name test-repo-1 --secret secret --vcs 1 --url https://github.com/Skarlso/test
9:51AM DBG Creating repository...
ID      NAME            URL                             VCS     CALLBACK-URL                                                    ATTACHED-COMMANDS     PROJECT-ID
1       test-repo-1     https://github.com/Skarlso/test 1       http://176.63.219.155:9998/rest/api/1/hooks/1/1/callback                              -1
```

Associate the test command to this repository with the following:

```
➜  krokctl git:(main) ✗ ./bin/darwin/amd64/krokctl relationship command add --command-id 1 --repository-id 1
Success!
```

And add the command association to the platform:

```
➜  krokctl git:(main) ✗ ./bin/darwin/amd64/krokctl relationship platform add --command-id 1 --platform-id 1
Success!
```

Now, if you navigate to the Github repository, you can keep sending it the ping event to test further functionality, like running commands and passing arguments correctly.

# Contributions

## Frontend

We would love to have help with Krok especially on the frontend side. None of us ( currently ) know frontend development too well.
We do plan on having a frontend for Krok alongside the CLI([krokctl](https://github.com/krok-o/krokctl)).

## Commands

Commands are never enough. Krok lives for the command it can execute on each webhook action. Without the commands Krok is almost useless.
Writing commands extends the universe of Krok and makes Krok itself more useful.

## Krok core

The server, which is this repository, can always use implementation of another platform for example. Or just general improvement of the code.