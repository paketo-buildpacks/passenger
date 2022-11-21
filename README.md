# Passenger Paketo Buildpack

## `gcr.io/paketo-buildpacks/passenger`

The Passenger CNB sets the start command for a given ruby application that runs on a Passenger server.

## Integration

This CNB writes a command, so there's currently no scenario we can
imagine that you would need to require it as dependency. If a user likes to
include some other functionality, it can be done independent of the Rake CNB
without requiring a dependency of it.

To package this buildpack for consumption:
```
$ ./scripts/package.sh
```
This builds the buildpack's source using GOOS=linux by default. You can supply another value as the first argument to package.sh.

## `buildpack.yml` Configurations

There are no extra configurations for this buildpack based on `buildpack.yml`.

## Compatibility

This buildpack is currently only supported on the Paketo Bionic and Jammy stack
distributions. A pre-compiled distribution of `curl` needed inside of the
buildpack is provided for the Paketo stacks (i.e.  `io.buildpacks.stack.jammy`
and `io.buildpacks.stacks.bionic`).
