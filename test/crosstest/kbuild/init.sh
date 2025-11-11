#!/bin/sh
mount -t proc proc /proc
mount -t sysfs sysfs /sys

setsid cttyhack sh
echo o > /proc/sysrq-trigger
