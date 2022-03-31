#!/usr/bin/env bash

set -euo pipefail

NAME="azapi2azurerm"
BUILD_ARTIFACT="${NAME}"
ARCHIVE_ARTIFACT="${NAME}_${VERSION}"
CHECKSUM_FILE="${ARCHIVE_ARTIFACT}_SHA256SUMS"
CHECKSUM_SIG="${ARCHIVE_ARTIFACT}_SHA256SUMS.sig"

OS_ARCH=("freebsd:amd64"
  "freebsd:386"
  "freebsd:arm"
  "freebsd:arm64"
  "windows:amd64"
  "windows:386"
  "linux:amd64"
  "linux:386"
  "linux:arm"
  "linux:arm64"
  "darwin:amd64"
  "darwin:arm64")


function checksum() {
  pwd
  ls
  for os_arch in "${OS_ARCH[@]}" ; do
    OS=${os_arch%%:*}
    ARCH=${os_arch#*:}
    EXT=$([ "$OS" == "windows" ] && echo ".exe" || echo "")
    echo "GOOS: ${OS}, GOARCH: ${ARCH}"
    (
      BIN_FILE="${BUILD_ARTIFACT}_${OS}_${ARCH}${EXT}"
      BIN_ZIP="${ARCHIVE_ARTIFACT}_${OS}_${ARCH}.zip"
      tar -a -c -f $BIN_ZIP $BIN_FILE
      rm -rf $BIN_FILE
      echo "$(certutil -hashfile $BIN_ZIP | head -2 | tail -1) $BIN_ZIP" >> $CHECKSUM_FILE
    )
  done
  cp $CHECKSUM_FILE $CHECKSUM_SIG
  cat $CHECKSUM_FILE
}

checksum


