To test locally:

```shell
# assume $output_dir is the output from the compilation step, with a tarball and a checksum in it

docker run -it \
  --volume $output_dir:/tmp/output_dir \
  --volume $PWD:/tmp/test \
  ubuntu:jammy \
  bash

# Now on the container

$ apt-get update && apt-get install ca-certificates -y

# Passing
$ /tmp/test/test.sh \
  --tarballPath /tmp/output_dir/curl_7.85.0_linux_jammy_c3a27edf.tgz \
  --expectedVersion 7.85.0
tarballPath=/tmp/output_dir/curl_7.85.0_linux_jammy_c3a27edf.tgz
expectedVersion=7.85.0
All tests passed!

# Failing
$ /tmp/test/test.sh \
  --tarballPath /tmp/output_dir/curl_7.85.0_linux_jammy_c3a27edf.tgz \
  --expectedVersion 7.84.0
tarballPath=/tmp/output_dir/curl_7.85.0_linux_jammy_c3a27edf.tgz
expectedVersion=7.84.0
Version 7.85.0 does not match expected version 7.84.0
```