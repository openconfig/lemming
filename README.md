# Lemming the Openconfig reference device implementation

## Purpose

To provide a reference implementation of a device which provides the collection
of customers APIs used for Openconfig. This includes:

* gNMI
* gNOI
* gRIBI
* P4RT
* BGP
* ISIS

to clearly and authoritatively specify the expected behavior of an
OpenConfig-compliant device, and to aid in its test development and
debugging. Anyone can use this reference implementation to quickly write device
tests independent of the availability of a real and compliant implementation.
The reference can also be used to augment the consumer contract when a
fake-derived test suite is delivered to network device vendors, serving as a
tool for consensus in the device-implementor <-> device-consumer relationship.

## Running the Fake gNMI Server

```bash
go run ./cmd/lemming --zapi_addr unix:/tmp/zserv.api --alsologtostderr
```

Wait for the message "lemming initialization complete".  
This might take a minute the first time to compile the large generated code.

Install gnmic: <https://gnmic.openconfig.net/>

```bash
// SetRequest configuring hostname
gnmic -a localhost:9339 --insecure -u foo -p bar --target fakedut set --update-path openconfig:/system/config/hostname --update-value rosesarered -e json_ietf

// SubscribeRequest/ONCE getting configured hostname
gnmic -a localhost:9339 --insecure -u foo -p bar --target fakedut subscribe --mode once --path openconfig:/system/config/hostname

// SubscribeRequest/ONCE getting hostname reflected as system state
gnmic -a localhost:9339 --insecure -u foo -p bar --target fakedut subscribe --mode once --path openconfig:/system/state/hostname
```

## Running integration tests

Prerequisites:

* [KNE](https://github.com/openconfig/kne) setup and cluster deployed

Deploy and Test:

* Optional: Build and load lemming container from source: `make load`
* Run integration tests: `go test ./integration_tests/...`


## Debugging Lemming

1. Load the debug image in the cluster: `make load-debug`
2. Modify a topology.pb.txt to start lemming with dlv.
```prototext
nodes: {
    name: "lemming"
    vendor: OPENCONFIG
    config: {
        command: "/dlv/dlv"
        args: "exec"
        args: "--headless"
        args: "--continue"
        args: "--accept-multiclient"
        args: "--listen=:56268"
        args: "--api-version=2"
        args: "/lemming/lemming"
        args: "--"
    }
}
```
3. Create the topology: `kne create <topofile>`.
4. Forward the debugger connection (this is blocking so run in seperate terminal): `kubectl port-forward -n <topo name> <node name> 56268:56268`
5. Attach to the debugger.
    1. Using VS Code: Run and Debug -> Connect to server
    2. Using dlv cli: `dlv connect localhost:56268`
        1. Required: Configure subsitute-path so dlv can resolve source code: `config substitute-path /build <abs path to lemming src>`

## Configuration

In general, Lemming should be configured through gNMI, however in some cases using flags is acceptable.

When to use gNMI:
* Config modelled in OC
* Values that may be modified at runtime

When to use flags:
* Startup options: (eg gNMI listen port)
* Immutable config
* Environment specific options (location of some resource)

Lemming uses Viper and pflags for configuration. Currently, only flags are supported (no env vars and no config file).
Flags are defined in cmd/lemming/lemming.go.


In KNE, flags are set using the `args` attribute in the topology file.
The Lemming operator also adds some mandatory flags to lemming for ease of use. These are flags that are always set
since they required to run lemming in a containerized environment.
They are defined at operator/controllers/lemming_controller.go.