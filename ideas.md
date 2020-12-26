# Krok Ideas and design decisions

## The URL

The hook url should be generated with an ID which will identify the hook itself.
This should be a UUID... But this also means that I need to pay attention to brute force
attacks. Something like, if there are enough misses, I ban the source IP for a while?

### Generating the ID

When a repository is added, a unique ID is generated for that repository. The hook
will have that ID. When the repository is created, we also ask where it is located, and
that will determine the Type of the hook. So when the hook is received through the `/hook/:id`
URL and the `id` is checked, we check if a repository with that ID exists. If so,
we get what the hook type is, and pass it to the appropriate provider.

## Plugins

The loader will be the one which will watch for new files being dropped into the folder.
It will do that in a separate go routine and store all commands in a cache map.

Need the map to be on the server?

*Update*: As to not having to update a multitude of SDKs, for the time being it is
decided to stick with Go plugins.

```go
type Plugin interface {
    Execute(payload string) (output string, outcome bool, err error)
}
```

This pretty much defines a plugin which can be executed.

*Update*: plugins will be uploaded as ZIP files to prevent filesystem dependency
so Krok can be ephemeral.

## Events

The hooks could follow an event system of some kind. I'm still debating on how to handle the
cron jobs. Run some kind of reconciliation loop? Some kind of timing? The crons could be
launched via the corn plugin. We should be able to stop and start them?

This means that on each tick it needs to check the settings.

## Plugins and Crons

A plugin first loaded, will be a blank plugin with no repositories and no cron schedule.
Once there is a cron schedule, you can stop/start the plugin, which will then perform it's
magic on all the assigned repositories.

Do I need a provider for the cron runner? The stop starter maybe, but I guess I deal with it
when I get there?

If there is a delete action in the watched folder, we mark that plugin as not working until the same
file with the same hash is placed into the folder. This will block the ability to update commands.
I guess that's fine, because it's a security feature. A command is immutable.

## Database

rel_repositories_command
rel_command_repositories

When a repository is deleted, delete the relationship from both sides. Same for a hook.

commands -> repositories
repositories -> commands

Relationship cannot be directly updated. It will be updated when a command is created and
assigned to a repository. It's not part of the command directly. It's stored separately.
This means it needs its own provider which will manage the relationship.

## Confidential data

Things like username and password for a repository should not be stored in the DB. I created
a Vault for that which will store confidential data for a repository with this format:
REPOID_USERNAME
REPOID_PASSWORD
REPOID_SSH_KEY

RepoID should be unique since it will be a uuid.
This data should be loaded on the fly when needed, and not stored with the repository object in memory.

## Commands

Don't forget to stream back the command's output.

## Connection

Remove the connection between repository and command. Figure out a better way to track the relationship. Maybe use the
same table?

Use proper foreign key definitions and cascading delete on the intermediary table which connects commands with repos.

Only a command can be assigned to a repository and not the other way around. So it makes sense to only have an
AddCommandToRepository because the other way around doesn't make sense.

## Users

Users will be done by logging in / registering via Google and OpenID. The platform will not provide registration and such.
The backend will use JWT to authenticate. The pure API login will be handled by API keys. The normal frontend will generate
the JWT using openid and google log in. The api will generate the JWT using the api key matching. Api Keys will be encrypted
and stored encrypted and matched encrypted.
