#!/bin/bash
set -euo pipefail

if [ $# -lt 2 ]; then
  echo "Usage: $0 <program> -- [args...]" >&2
  exit 1
fi

program="$1"
shift

if [ "$1" != "--" ]; then
  echo "Expected '--' separator" >&2
  exit 1
fi
shift

args_file=$(mktemp)
stdout_pipe=$(mktemp -u)
stderr_pipe=$(mktemp -u)
control_pipe=$(mktemp -u)

cleanup() {
  rm -f "$args_file"
  rm -f "${stdout_pipe}.in" "${stdout_pipe}.out"
  rm -f "${stderr_pipe}.in" "${stderr_pipe}.out"
  rm -f "${control_pipe}.in" "${control_pipe}.out"
}
trap cleanup EXIT

printf '%s\0' "$@" > "$args_file"

mkfifo "${stdout_pipe}.in"
mkfifo "${stdout_pipe}.out"
mkfifo "${stderr_pipe}.in"
mkfifo "${stderr_pipe}.out"
mkfifo "${control_pipe}.in"
mkfifo "${control_pipe}.out"

cat "${stdout_pipe}.out" &
stdout_cat_pid=$!

cat "${stderr_pipe}.out" >&2 &
stderr_cat_pid=$!

qemu-system-aarch64 \
  -machine virt,gic-version=3 \
  -cpu host \
  -accel hvf \
  -nographic \
  -m 1G \
  -kernel ./kbuild/out/arm64/Image \
  -append "console=ttyAMA0 loglevel=0" \
  -fw_cfg name=opt/program,file="$program" \
  -fw_cfg name=opt/args,file="$args_file" \
  -device virtio-serial-device \
  -chardev pipe,id=stdout,path="$stdout_pipe" \
  -device virtserialport,chardev=stdout,nr=2 \
  -chardev pipe,id=stderr,path="$stderr_pipe" \
  -device virtserialport,chardev=stderr,nr=3 \
  -chardev pipe,id=control,path="$control_pipe" \
  -device virtserialport,chardev=control,nr=4 \
  >/dev/null 2>&1 &

qemu_pid=$!

read exit_code < "${control_pipe}.out"

wait $stdout_cat_pid 2>/dev/null || true
wait $stderr_cat_pid 2>/dev/null || true
wait $qemu_pid 2>/dev/null || true

exit $exit_code
