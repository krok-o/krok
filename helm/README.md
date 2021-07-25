# Storage

If `storage` is enabled use mount points defined in the storage section.
Otherwise, use in container storage, aka. `/tmp/krok/*`.

# Dependencies

In order for Krok to run, first, some configMaps and a secret needs to run and the database.
We define these as dependencies.

```
helm dependency list helm
NAME            VERSION REPOSITORY      STATUS
krok-config     0.1.0                   unpacked
krok-database   0.1.0                   unpacked
```

Initialise and build the dependencies first, then install Krok itself.