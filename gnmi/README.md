# gNMI Service Reference Implementation (central datastore)

`gnmi.GNMIServer` serves as the reference implementation of gNMI for lemming and the central datastore to which reconcilers can be attached.

## External Clients

External clients use gNMI to configure and read state of the reference device.
## Internal Reconcilers

Lemming is designed as an eventually consistent system where the gNMI cache is used as the central datastore. A reconciler is a set of initialization methods and long-running methods (either of which are optional) that monitor the intended configuration
(or other internal state) and modify the operational state of the device to reach the intended state. Reconcilers are attached to the gNMI cache and  are responsible reading the intended config and reconciling state.
Reconcilers can optionally validate incoming SetRequest to prevent semantically incorrect values from being applied. For example, the interface reconciler can validate that MTU > 64 for Ethernet interfaces.

Reconcilers are creating by implementing the [Reconciler interface](https://pkg.go.dev/github.com/openconfig/lemming/gnmi/reconciler#Reconciler), and optionally using the a reconciler Builder to simplify the creation.

See the [fakedevice package](fakedevice/fakedevice.go) for examples of reconcilers.

### Reconciliation

Reconcilers must implement the the reconciler.Reconciler interface. The reconciler interface is intentionally minimum to minimize the burden of creating reconcilers and allow maximum flexibility when creating a reconciler.

#### Reacting to Values

Reconcilers can react to values by either calling a raw gNMI operation using the provided GNMIClient argument or by starting a Watch using ygnmi (optionally using WithWatch). ygnmi.Watch can be used monitoring the current state of a set of OpenConfig subtrees/leaves using gNMI.Subscribe.

#### Updating Values

Reconcilers can update values in the cache by send raw gNMI set requests using the client or using ygnmi. The gnmiclient package contains special helpers that allow using ygnmi.Update/Replace/Delete on state paths. Generally, internal reconcilers should not be updating config paths as those are set by external users.

NOTE: Currently, Set requests are required to be schema-compliant.
NOTE: Replace/Update on non-leaf state values is not yet. See [issue](https://github.com/openconfig/lemming/issues/67).

### Validation

Reconcilers can validate intended config before it is applied to the cache. Reconcilers must return which path prefixes they are interested in validating. For any SetRequest that modifies those paths, the Validate function is called on the global intended config.

Note: this currently implemented as any reconciler that returns > 0 paths is called on every SetRequest.

### Reconciler Builders

To simplify the creation of reconcilers, the reconciler package contains a builder API.
The reconciler package provides a builder API for creating reconcilers with common use cases:

```go
rec := reconciler.NewTypedBuilder[string]("hostname").
    WithWatch(hostnamePath.Config(), func(ctx context.Context, c *ygnmi.Client, v *ygnmi.Value[string]) error { // Watch hostname config: set hostname, and update state.
        hostname, ok := v.Val()
        if !ok {
            resetHostName()
        }
        setHostName(hostName)
        if _, err := gnmiclient.Update(ctx, c, hostnamePath.State(),hostname); err != nil {
            return log.Warningf("Failed to set hostname: %v",err)
        }
        return ygnmi.Continue
    }).Build()
```

Alternatively, a custom reconciler implementation matching the Reconciler interface can be created and registered on the gnmi server.
