# Lemming CI

## Presubmit

Any PR triggers the presubmit job. The presubmit creates GCE VM, sets up KNE instance, and runs tests in integration_tests folder.

See [remote-builder](https://github.com/GoogleCloudPlatform/cloud-builders-community/tree/master/remote-builder) for details on the remote builder.

## Release

### Prerequisites

Ensure your Application Default Credentials (ADC) quota project is set to `openconfig-lemming` so that Cloud Build API calls are authorized and billed to the correct project:

```bash
gcloud auth application-default set-quota-project openconfig-lemming
```

### Lemming

To release a new version of lemming run: `go run ./cmd/lemctl internal release lemming VERSION`

The following steps are run:  
1. Trigger prerelease tests in Cloud Build: lemming integration tests (see lemming-test.sh).
    1. Runs using the latest released version of the operator and builds lemming from HEAD.
2. Create and push git tag.
3. Push lemming to GAR.

### Operator

To release a new version of operator run: `go run ./cmd/lemctl internal release operator VERSION`

The following steps are run:  
1. Trigger prerelease tests in Cloud Build: deploys a simple 2 lemming topology (see operator-test.sh).
    1. Runs using the latest released version of lemming and builds operator from HEAD.
2. Create and push git tag.
3. Push operator to GAR.

### After push (not automated)

1. Modify `operator/config/manager/kustomization.yaml` `newTag` to the new version.
2. Run `kubectl kustomize operator/config/default > <kne-repo>/manifests/controllers/lemming/manifest.yaml`
3. Create PR to update manifest in kne.