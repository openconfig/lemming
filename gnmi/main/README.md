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

To clearly and authoritatively specify the expected behavior of Google’s needs
for a device, and to aid in its test development and debugging. Anyone can use
this reference implementation to quickly write device tests independent of the
availability of a real and compliant implementation. The reference can also be
used to augment the consumer contract when a fake-derived test suite is
delivered to network device vendors, serving as a tool for consensus in the
vendor-operator relationship.

## Running the Fake gNMI Server

```bash
go run main.go
```

```bash
gnmic -a localhost:1234 --insecure subscribe --mode stream --path openconfig:/system/state/current-datetime -u foo -p bar --target fakedut
```

For running on KNE (also experimental), see
[here](https://github.com/wenovus/ondatra/tree/fake-prototype-0/fakebind)
