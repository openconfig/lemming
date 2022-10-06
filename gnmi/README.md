# gNMI Service Reference Implementation (central datastore)

`gnmi.GNMIServer` serves as the reference implementation of gNMI for lemming and the central datastore to which reconcilers can be attached.

## External Clients

External clients use gNMI to configure and read state of the reference device.
## Internal Reconcilers

Lemming is designed as an eventually consistent system where the gNMI cache is used as the central datastore. Reconcilers are attached to the gNMI cache and  are responsible reading the intended config and reconciling state.
Reconcilers can optionally validate incoming SetRequest to prevent semantically incorrect values from being applied. For example, the interface reconciler can validate that MTU > 64 for Ethernet interfaces.

Reconcilers are creating by implementing the Reconciler interface, and optionally using the a reconciler Builder to simplify the creation.

### Reconciliation

Reconcilers must implement the the reconciler.Reconciler interface. The reconciler interface is intentionally minimum to minimize the burden of creating reconcilers and allow maximum flexibility when creating a reconciler.

#### Reacting to Values

Reconcilers can react to values by either creating a "raw" gNMI Subscription using the client or by starting a Watch using ygnmi (optionally using WithWatch). The values are modelled in OpenConfig, using Watch is generally simpler and easier to write.

#### Updating Values

Reconcilers cab update values in the cache by send raw gNMI set requests using the client or using ygnmi. The gnmiclient package contains special helpers that allow using ygnmi.Update/Replace/Delete on state paths. Generally, internal reconcilers should not be updating config paths as those are set by external users.

NOTE: Currently, Set requests are required to be schema-compliant.

### Validation

Reconcilers validate intended config before it is applied to the cache. Reconcilers must return that path prefixes they are interested in validating. For any SetRequest that modifies those paths, the Validate function is called on the global intended config.

Note: this currently implemented as any reconciler that returns > 0 paths is called on every SetRequest.

### Reconciler Builders

To simplify the creation of reconcilers, the reconciler package contains a builder API. The builder API contains funcs for common use cases (like creating a reconciler that runs a ygnmi Watch). If the builder API doesn't fit the reconciler, it's always possible to implement the reconciler interface directly.
