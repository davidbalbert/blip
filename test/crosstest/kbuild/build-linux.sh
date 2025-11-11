#!/bin/sh
set -euo pipefail

export KBUILD_BUILD_TIMESTAMP="${BUSYBOX_DATE}"
export KBUILD_BUILD_USER=root
export KBUILD_BUILD_HOST=build


mkdir -p /src
tar -C /src -xf "/pkg/linux-${LINUX_VERSION}.tar.xz"
mkdir -p "/build/linux-${LINUX_VERSION}"
cd "/build/linux-${LINUX_VERSION}"
cp /config/linux.config .config
make -C "/src/linux-${LINUX_VERSION}" O="/build/linux-${LINUX_VERSION}" ARCH=arm64 olddefconfig
make -j$(nproc) -C "/src/linux-${LINUX_VERSION}" O="/build/linux-${LINUX_VERSION}" ARCH=arm64

cp /build/linux-${LINUX_VERSION}/arch/arm64/boot/Image /out
