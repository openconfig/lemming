# Fake gNMI Service (central datastore)

`gnmit.GNMIServer` serves as a central datastore to which subscribing
"tasks/agents" can be attached.

There are two ways to subscribe to the central datastore:

## 1 Subscribing using a `gnmit.Task` interface

To use this method, define a function that satisfies the `gnmit.TaskRoutine`
function signature, and register the task in [lemming.go](../lemming.go).

See [lemming.go](../lemming.go) to see how
[testagentlocal](testagentlocal/interface.go), an example task implementation,
is registered.

## 2 Subscribing using ygnmi

This method uses [`ygnmi`](https://github.com/openconfig/ygnmi) to receive
updates from the central database.

See [testagent](testagent/testagent.go) for an example.
