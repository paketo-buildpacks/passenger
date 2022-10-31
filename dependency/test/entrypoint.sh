#!/bin/bash

set -euo pipefail
shopt -s inherit_errexit

extract_tarball() {
  rm -rf curl
  mkdir curl
  tar --extract \
    --file "$1" \
    --directory curl
}

set_ld_library_path() {
  export LD_LIBRARY_PATH="$PWD/curl/lib:${LD_LIBRARY_PATH:-}"
}

check_version() {
  expected_version=$1
  actual_version="$(./curl/bin/curl -V | head -n1 | awk '{ print $2 }')"
  if [[ "${actual_version}" != "${expected_version}" ]]; then
    echo "Version ${actual_version} does not match expected version ${expected_version}"
    exit 1
  fi
}

check_server() {
  output="$(mktemp)"
  if ! ./curl/bin/curl -fsS https://example.org > "${output}"; then
    cat "${output}"
    exit 1
  fi
}

main() {
  local tarballPath expectedVersion
  tarballPath=""
  expectedVersion=""

  while [ "${#}" != 0 ]; do
    case "${1}" in
      --tarballPath)
        tarballPath="${2}"
        shift 2
        ;;

      --expectedVersion)
        expectedVersion="${2}"
        shift 2
        ;;

      "")
        shift
        ;;

      *)
        echo "unknown argument \"${1}\""
        exit 1
    esac
  done

  if [[ "${tarballPath}" == "" ]]; then
    echo "--tarballPath is required"
    exit 1
  fi

  if [[ "${expectedVersion}" == "" ]]; then
    echo "--expectedVersion is required"
    exit 1
  fi

  echo "tarballPath=${tarballPath}"
  echo "expectedVersion=${expectedVersion}"

  extract_tarball "${tarballPath}"
  set_ld_library_path
  check_version "${expectedVersion}"
  check_server

  echo "All tests passed!"
}

main "$@"
