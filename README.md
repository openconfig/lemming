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
This might take a minute to compile the first time due to the size of the generated OpenConfig code.

Install gnmic: <https://gnmic.openconfig.net/>

```bash
// SetRequest configuring hostname
gnmic -a localhost:6030 --insecure -u foo -p bar --target fakedut set --update-path openconfig:/system/config/hostname --update-value rosesarered -e json_ietf

// SubscribeRequest/ONCE getting configured hostname
gnmic -a localhost:6030 --insecure -u foo -p bar --target fakedut subscribe --mode once --path openconfig:/system/config/hostname

// SubscribeRequest/ONCE getting hostname reflected as system state
gnmic -a localhost:6030 --insecure -u foo -p bar --target fakedut subscribe --mode once --path openconfig:/system/state/hostname
```

## Running integration tests

Prerequisites:

* [KNE](https://github.com/openconfig/kne) setup and cluster deployed

Setup:

* Deploy the operator: `kubectl apply -k operator/config/default`
* Optional: Build and load lemming container from source: `make load`

Deploy and Test:

* Run integration tests: `go test ./integration_tests/...`
