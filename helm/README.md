# TODO

- [x] Configurable ingress
- [x] Configurable mount paths
- [x] Secrets and ConfigMaps if needed
- [x] Mount DB secret into krok server and use that

Extract Krok and the Database as a separate Chart and use Dependencies in order to define them.
Update the database url to be a service url defined by Kubernetes and the database service name.
(Dependencies)[https://helm.sh/docs/topics/charts/] -- follow this thing describe here and create
sub-charts.

# Storage

If `storage` is enabled use mount points defined in the storage section.
Otherwise, use in container storage, aka. `/tmp/krok/*`.