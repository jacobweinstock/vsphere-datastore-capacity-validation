# vSphere Validator

[![Test and Build](https://github.com/jacobweinstock/vvalidator/workflows/Test%20and%20Build/badge.svg)](https://github.com/jacobweinstock/vvalidator/actions?query=workflow%3A%22Test+and+Build%22)
[![Go Report](https://goreportcard.com/badge/github.com/jacobweinstock/vvalidator)](https://goreportcard.com/report/github.com/jacobweinstock/vvalidator)

Generic vSphere Validations

## Usage

### Build

Build the binary by running `make build`

### Commands

All commands return a response to stdout, and a file (defaults to `./response.json`) in json format.

#### Capacity

Calculates if a datastore has space for the requested disk(s) size.  
The disk(s) size can be collected from a VM if the name is give to `--vmName` (flag) or `VVALIDATOR_VMNAME` (env var).

```bash
# validate a datastore has capacity for a VM
vvalidator capacity \
  --vmName bionic-server-cloudimg-amd64 \
  --url 10.96.160.151 \
  --datacenter Datacenter-01 \
  --datastore Datastore-01 \
  --user administrator@vsphere.local \
  --password 'secret'
```


```bash
# Environment variables can be used in place of the CLI flags
VVALIDATOR_USER=myuser \
VVALIDATOR_PASSWORD=secret \
VVALIDATOR_URL=vcenter.example.org \
VVALIDATOR_DATACENTER=dc1 \
VVALIDATOR_DATASTORE=ds1 \
VVALIDATOR_VMNAME=mytemplate \
vvalidator capacity
```

##### Response Object

For more details on the data types, the go `capacityResponse` struct can be found here: `cmd/response.go`
```json
{
  "errorMsg": "",
  "freeSpaceInGBs": 0,
  "level": "",
  "msg": "",
  "requestedSpaceInGBs": 0,
  "responseFile": "",
  "spaceAvailable": false,
  "success": false,
  "time": ""
}
```

#### Size

Returns the size of all disks, in GBs, attached to a VM or template. 

```bash
vvalidator size \
  --vmName bionic-server-cloudimg-amd64 \
  --url 10.96.160.151 \
  --datacenter Datacenter-01 \
  --user administrator@vsphere.local \
  --password 'secret'
``` 

Environment variables can be used in place of the CLI flags
```bash
VVALIDATOR_USER=myuser \
VVALIDATOR_PASSWORD=secret \
VVALIDATOR_URL=vcenter.example.org \
VVALIDATOR_DATACENTER=dc1 \
VVALIDATOR_VMNAME=mytemplate \
vvalidator size
```

##### Response Object

For more details on the data types, the go `sizeResponse` struct can be found here: `cmd/response.go`
```json
{
  "errorMsg": "",
  "level": "",
  "msg": "",
  "responseFile": "",
  "success": false,
  "time": "",
  "totalDiskSize": 0,
  "vmName": ""
}
```

### Container Image

A container image is available at `docker pull ghcr.io/jacobweinstock/vvalidator:latest`  

#### USAGE

Build the image locally with `make build`

```bash
# CAPACITY using cli flags
docker run -it --rm ghcr.io/jacobweinstock/vvalidator size \
  --datacenter DC01 \
  --password 'secret' \
  --user user \
  --url 10.96.160.151 \
  --vmName bionic-server-cloudimg-amd64

# using env vars
docker run -it --rm \
  -e VVALIDATOR_URL=10.96.160.151 \
  -e VVALIDATOR_DATACENTER=Datacenter-01 \
  -e VVALIDATOR_DATASTORE=Datastore-01 \
  -e VVALIDATOR_USER=admin \
  -e VVALIDATOR_PASSWORD='secret' \
  -e VVALIDATOR_VMNAME=bionic-server-cloudimg-amd64 \
  -v ${PWD}/response.json:/response.json \
  ghcr.io/jacobweinstock/vvalidator capacity


# SIZE using cli flags
docker run -it --rm ghcr.io/jacobweinstock/vvalidator size \
  --datacenter DC01 \
  --password 'secret' \
  --user user \
  --url 10.96.160.151 \
  --vmName bionic-server-cloudimg-amd64

# using env vars
docker run -it --rm \
  -e VVALIDATOR_URL=10.96.160.151 \
  -e VVALIDATOR_DATACENTER=Datacenter-01 \
  -e VVALIDATOR_USER=admin \
  -e VVALIDATOR_PASSWORD='secret' \
  -e VVALIDATOR_VMNAME=bionic-server-cloudimg-amd64 \
  -v ${PWD}/response.json:/response.json \
  ghcr.io/jacobweinstock/vvalidator size
```