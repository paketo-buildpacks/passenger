To test locally:

```shell
# Assume $output_dir is the output from the compilation step, with a tarball and a checksum in it.
# Note that the wildcard is not quoted, to allow globbing

# Passing
$ ./test.sh \
  --tarballPath ${output_dir}/*.tgz \
  --expectedVersion 7.85.0
tarballPath=/tmp/output_dir/curl_7.85.0_linux_bionic_f0e0bef7.tgz
expectedVersion=7.85.0
All tests passed!

# Failing
$ ./test.sh \
  --tarballPath ${output_dir}/*.tgz \
  --expectedVersion 999.999.999
tarballPath=/tmp/output_dir/curl_7.85.0_linux_bionic_f0e0bef7.tgz
expectedVersion=999.999.999
Version 7.85.0 does not match expected version 999.999.999
```