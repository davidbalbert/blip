#!/bin/sh
mount -t proc proc /proc
mount -t sysfs sysfs /sys

cat /sys/firmware/qemu_fw_cfg/by_name/opt/program/raw > /tmp/program
chmod +x /tmp/program

cat /sys/firmware/qemu_fw_cfg/by_name/opt/args/raw > /tmp/args

# Build the command with arguments
set --
while IFS= read -r -d '' arg; do
  set -- "$@" "$arg"
done < /tmp/args

/tmp/program "$@" </dev/vport0p1 >/dev/vport0p2 2>/dev/vport0p3
exit_code=$?

echo $exit_code > /dev/vport0p4

echo o > /proc/sysrq-trigger
