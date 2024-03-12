# How to add a new feature

This section describes how to add a new feature to the dataplane using a hypothetical and very contrived example (and just about the worst case).

Let's implement a new feature: based on the hypothetical `foo` protocol packet header, do a shortest prefix match on dst ip and on match randomize the payload bytes.

## Forwarding support for feature

The first step is to check if the forwarding supports everything we need. The answer is usually yes, but not in this case.

1. There is no shortest prefix forwarding table, so let's add it.
   1. Add a new value to the `TableType` enum in the [forwarding_table](../../proto/forwarding/forwarding_table.proto).
   2. Add a new entry desc proto message for the new table and it to the `oneof entry` in `EntryDesc`.
   3. Implement the table in a new package in the [fwdtable](../forwarding/fwdtable/).
   4. Tables must implement the `fwdtable.Table` interface and register a builder (`fwdtable.Builder`). (see other packages for examples).
2. There is no randomized payload action, so let's add it.
   1. Add a new value to the `ActionType` enum in the [forwarding_action](../../proto/forwarding/forwarding_action.proto).
   2. Add a new action desc proto message for the new action and it to the `oneof action` in `ActionDesc`.
   3. Implement the action in a new file in the [fwdaction](../forwarding/fwdactions/actions).
   4. Tables must implement the `fwdaction.Action` interface and register a builder (`fwdaction.Builder`). (see other packages for examples).
3. There is no `foo` header parser, so let's add it.
   1. Add a new enum value to the PacketFieldNum [forwarding_common](../../proto/forwarding/forwarding_common.proto).
   2. Update the forwarding/protocol package with the logic to parse this new header.

## saiserver implementation

Now that the datataplane supports all the features we need, we have to provide an API to configure it. In the saiserver package, there are structs that
embed the unimplemented server for all the saipb services.

For this example, we are going to assume that this feature is defined in the saipb API as follows.

```go
func (f *foo) CreateFoo(ctx context.Context, req *saipb.CreateFooRequest) (*saipb.CreateFooResponse, error) {
}

func (f *foo) RemoveFoo(ctx context.Context, req *saipb.RemoveFooRequest) (*saipb.RemoveFooResponse, error) {
}
```

Let's assume we want to implement the CreateFoo and RemoveFoo RPCs.

1. First, we need to set up the table and add it as a step in the forwarding pipeline.
   1. We need to create the table. Since there's only one table, creating it in the CreateSwitch RPC is usually the right place.
      1. In the CreateSwitch, create a new shortest prefix table table, which we can call the "foo-table"
   2. Now, we need to determine where in the pipeline should this action take place:
      1. For this example, let's pretend it belongs after updating the DST MAC from the neighbor table.
      2. Most of the pipeline is static and set using input actions on every port. If we look at [ports.go](../saiserver/ports.go), we can see all the actions applied on every port.
      3. In this, we are looking up what to put in our "foo-table", so we add a new lookup action to "foo-table" to slices of input actions.
2. Next, we need to implement the CreateFoo and RemoveFoo RPC.
   1. In CreateFoo, we need convert the CreateFooRequest to a `fwdpb.TableEntryAddRequest`
      1. Implementations will vary, there are plenty of examples in the saiserver package.
      2. Make sure that we add randomization action to the entry.
   2. In RemoveFoo, we need to do mostly the same.

>Note: A new feature may just require supporting an additional attribute in already implemented. If the feature is not in SAI API, then it can be configured using the forwarding API directly or by patching the generated protos.

## gNMI reconcilation

Not done yet! Now that the feature is available via the saiserver API, we may also want to configure it using gNMI.

1. In dplanerc package, let's add a new [reconciler](../../gnmi/reconciler/reconciler.go).
   1. First, we need to determine which OpenConfig paths reference this feature (if any). Make sure the relevant YANG file is added in [generate script](../../gnmi/generate.sh).
      1. The [OpenConfig site](https://openconfig.net/projects/models/paths/index.html) is the best reference.
   2. gNMI and Lemming use an eventual consistency, pub-sub model, so in the reconciler we need to subscribe to the config paths for our feature.
      1. Generally, using ygnmi.Watch and a batch subscription is useful.
   3. If the config does not equal the state, we need to reconcile the differences.
      1. This usually looks like, if config != state, call saiserver RPC to fix (Create or Remove)
      2. If the call was successful, then update the state.
      3. Note: Lemming has a special `gnmiclient` client that allows ygnmi.Set on state paths in order to have internal state managers/agents reflect the correct operational state of the virtual device.
      4. Note: For performance reasons: it is best to use BatchUpdate (not replace) individual leaves.

Almost there! Now write tests.
