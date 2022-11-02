# Lemming CI

## Presubmit

Any PR triggers the presubmit job. The presubmit creates GCE VM, sets up KNE instance, and runs tests in integration_tests folder.

See [remote-builder](https://github.com/GoogleCloudPlatform/cloud-builders-community/tree/master/remote-builder) for details on the remote builder.