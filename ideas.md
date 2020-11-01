# Krok Ideas and design decisions

## The URL

The hook url should be generated with an ID which will identify the hook itself.
This should be a UUID... But this also means that I need to pay attention to brute force
attacks. Something like, if there are enough misses, I ban the source IP for a while?

## Plugins

The loader will be the one which will watch for new files being dropped into the folder.
It will do that in a separate go routine and store all commands in a cache map.

Need the map to be on the server?

*Update*: So as to not having to update a multitude of SDKs, for the time being it is
decided to stick with Go plugins.

```go
type Plugin interface {
    Execute(payload string) (output string, outcome bool, err error)
}
```

This pretty much defines a plugin which can be executed.

## Events

The hooks could follow an event system of some kind. I'm still debating on how to handle the
cron jobs. Run some kind of reconciliation loop? Some kind of timing? The crons could be
launched via the corn plugin. But we should be able to stop and start them?

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
