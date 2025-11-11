#!/bin/sh
set -euo pipefail

export SOURCE_DATE_EPOCH=$(date -d "${BUSYBOX_DATE}" +%s)
export KCONFIG_NOTIMESTAMP=1

mkdir -p /src
tar -C /src -xf "/pkg/busybox-${BUSYBOX_VERSION}.tar.bz2"
mkdir -p "/build/busybox-${BUSYBOX_VERSION}"
cd "/build/busybox-${BUSYBOX_VERSION}"
cp /config/busybox.config .config
make -j$(nproc) -C "/src/busybox-${BUSYBOX_VERSION}" O="/build/busybox-${BUSYBOX_VERSION}" ARCH=arm64 CROSS_COMPILE=aarch64-linux-gnu- CONFIG_PREFIX=/build/rootfs
make -j$(nproc) CONFIG_PREFIX=/build/rootfs install
