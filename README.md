dgbrx
=====
> Dgraph Backup and Restore X

`dgbrx` is a Go commandline tool which helps to do a backup, restore
or clean on a Dgraph Cloud (aka slash / managed) instance.

Installation ✨
---------------
```bash
cd cmd/dgbrx
go build .
./dgbrx --help
```

Usage 🤔 
--------
There are primarily three commands for `dgbrx`: 
`backup`, `restore` and `clean`. `dgbrx backup` requests a backup 
from the provided Dgraph Instance, and waits for the backup to 
complete, which is then written to disk. `dgraph restore` can
restore the backups to a dgraph instance, and `dgraph clean` drops
the schema and all the data along with it.

Workflow 🔧
-----------
A general backup-restore workflow for dgraph is given below:
```bash
dgbrx backup --url https://some-cool-url.region.gcp.cloud.dgraph.io/admin \
  --api-key "SUPERSECRETAPIKEY"
  
dgbrx restore --url https://another-cool-url.region.gcp.gcloud.dgraph.io/admin \
  --api-key "SUPERSECRETAPIKEYBUTDIFFERENTONE" \
  --json g01.json.gz \
  --schema g01.schema.gz
```

Contributing 🔍
---------------
Make sure you adhere to Go formatting guidelines when contributing 
to this repository
```bash
go fmt 
```

Roadmap 🛣️
---------
- [ ] Add support for multiple storage backends (Google Cloud Storage, S3 Bucket, etc.)
- [ ] Implement unmanaged Dgraph Instance backup and restore.

Motivation 💪
-------------
[dgbr](https://github.com/AugustDev/dgbr), another open source 
dgraph backup and restore software did not support cloud.dgraph.io
(managed Dgraph Instances), since they use a slightly different API.

License ⚖
----------
This software is licensed under the [GNU Lesser General Public License v3](./LICENSE).
