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
go run cmd/lemming/lemming.go
```

```bash
gnmic -a localhost:6030 --insecure subscribe --mode stream --path openconfig:/system/state/current-datetime -u foo -p bar --target fakedut
```

For running on KNE (also experimental), see
[here](https://github.com/wenovus/ondatra/tree/fake-prototype-0/fakebind)
