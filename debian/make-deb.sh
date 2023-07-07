#!/bin/bash
set -e

APP_NAME="pi-hole-manager"
WORKDIR="/tmp/${APP_NAME}"

#set -x

test -d "${WORKDIR}" || mkdir -p "${WORKDIR}"
cp -ar /src "${WORKDIR}"
cd "${WORKDIR}/src"
debuild -us -uc
test -d /src/debian/build || mkdir /src/debian/build
cp ../${APP_NAME}* /src/debian/build/
echo -e "\nBuild artifacts are available on debian/build/ directory:"
ls -1 /src/debian/build/
