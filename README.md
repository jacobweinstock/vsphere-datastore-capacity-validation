# vSphere Validator

Generic vSphere Validations

## Usage

### Capacity

Calculates if a datastore has space for the requested disk(s) size. The disk(s) size can be collected from a VM if the name is give to `--name` (flag) or `VVALIDATOR_VMNAME` (env var). A response is return to stdout and a file (defaults to `./response.json`) in json format. 

validate a datastore has capacity for a VM
```bash
vvalidator capacity --vmName bionic-server-cloudimg-amd64 --url 10.96.160.151 --datacenter Datacenter-01 --datastore Datastore-01 --user administrator@vsphere.local --password 'secret'
```

Environment variables can be used in place of the CLI flags
```bash
 VVALIDATOR_USER=myuser VVALIDATOR_PASSWORD=secret VVALIDATOR_URL=vcenter.example.org VVALIDATOR_DATACENTER=dc1 VVALIDATOR_DATASTORE=ds1 VVALIDATOR_VMNAME=mytemplate vvalidator capacity
```

### Size

Returns the size of all disks, in GBs, attached to a VM or template. A response is return to stdout and a file (defaults to `./response.json`) in json format. 

```bash
vvalidator size --vmName bionic-server-cloudimg-amd64 --url 10.96.160.151 --datacenter Datacenter-01 --user administrator@vsphere.local --password 'secret'
``` 

Environment variables can be used in place of the CLI flags
```bash
 VVALIDATOR_USER=myuser VVALIDATOR_PASSWORD=secret VVALIDATOR_URL=vcenter.example.org VVALIDATOR_DATACENTER=dc1 VVALIDATOR_VMNAME=mytemplate vvalidator size
```