# TODO

- [ ] Configurable ingress
- [ ] Configurable mount paths
- [ ] Secrets and ConfigMaps if needed
- [ ] Mount DB secret into krok server and use that

# Storage

If `storage` is enabled use mount points defined in the storage section.
Otherwise, use in container storage, aka. `/tmp/krok/*`.