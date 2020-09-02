# vSphere Validator

Generic vSphere Validations

## Usage

* validate a datastore has capacity for a VM
```bash
vvalidator capacity --vm bionic-server-cloudimg-amd64 --url 10.96.160.151 --datacenter Datacenter-01 --datastore Datastore-01 --user administrator@vsphere.local --password 'secret'
```
* get size (in GBs) of all the disk(s) attached to a VM
```bash
vvalidator size --vm bionic-server-cloudimg-amd64 --url 10.96.160.151 --datacenter Datacenter-01 --user administrator@vsphere.local --password 'secret'
```