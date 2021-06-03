[![GitHub version](https://badge.fury.io/gh/mjuz-iac%2Fpulumi.svg)](https://badge.fury.io/gh/mjuz-iac%2Fpulumi)

# Pulumi for µs

This is a fork of [**Pulumi's Infrastructure as Code SDK**](https://github.com/pulumi/pulumi)
for [µs Infrastructure as Code](https://mjuz.rocks).
It adds pruning of resources that depend on an unsatisfied µs wish to the deployment engine.
This prevents the deployment of resources that (transitively) depend on an unsatisfied wish and
triggers their deletion if they were deployed before.
µs requires the Pulumi CLI of this fork instead of the official Pulumi CLI for its correct operation.

## Installation

This project is a drop-in replacement for the official Pulumi CLI.
For development and building from source, it requires the same setup as Pulumi,
documented in [CONTRIBUTING.md](CONTRIBUTING.md).

To build and install the CLI, run `make install` in the repository's root directory.
It will download all dependencies, build the CLI and install to `/opt/pulumi`.
Make sure that a `/opt/pulumi` directory exists on your system before and that your user has the required read, write,
and execute permissions for it.

µs programs will invoke the CLI as `pulumi`.
Prepend `/opt/pulumi:/opt/pulumi/bin` to the `PATH` environment variable to ensure this CLI version is used instead
of another Pulumi installation on your system.
To suppress the upgrade version warning at start, set the environment variable `PULUMI_SKIP_UPDATE_CHECK`.
This fork uses its own versioning scheme that will always result in showing an upgrade warning otherwise.

To check whether your setup is correct, run `pulumi version`. It should only print the version number,
which for this fork will always contain `-mjuz`, e.g., `v1.0.0-mjuz` or `v1.0.0-mjuz+dirty`.
If `-mjuz` is not in the printed version number, you are most likely executing a wrong Pulumi CLI installation.

## Docker

We provide a docker image with of this CLI on Docker Hub as [mjuz/pulumi](https://hub.docker.com/r/mjuz/pulumi).
It contains all dependencies and this repository's content in `/var/pulumi`,
from which the µs Pulumi CLI is already installed.
It can be re-installed inside the container by running `make install` in `/var/pulumi`.
