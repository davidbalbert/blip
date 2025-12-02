#!/bin/bash
set -euo pipefail

if [ $# -lt 1 ]; then
  echo "Usage: $0 <program> [args...]" >&2
  exit 1
fi

program="$1"
shift

tmpdir=$(mktemp -d)
cleanup() {
  local pids
  pids=$(jobs -p)
  if [ -n "$pids" ]; then
    kill $pids 2>/dev/null || true
  fi
  rm -rf "$tmpdir"
}
trap cleanup EXIT

for p in stdin stdout stderr control; do
  mkfifo "$tmpdir/$p.in" "$tmpdir/$p.out"
done

# Build payload: [8-byte total_size][8-byte program_size][program][8-byte args_size][args]
send_payload() {
  local program_size args_data args_size total_size
  program_size=$(stat -f '%z' "$program")
  args_data=$(printf '%s\0' "$@")
  args_size=${#args_data}
  total_size=$((8 + program_size + 8 + args_size))

  printf '%08d%08d' "$total_size" "$program_size"
  cat "$program"
  printf '%08d%s' "$args_size" "$args_data"
}

cat > "$tmpdir/stdin.in" &
cat "$tmpdir/stdout.out" &
cat "$tmpdir/stderr.out" >&2 &

qemu-system-aarch64 \
  -machine virt,gic-version=3 -cpu host -accel hvf \
  -nographic -m 1G \
  -kernel ./kbuild/out/arm64/Image \
  -append "console=ttyAMA0 loglevel=0" \
  -device virtio-serial-device \
  -chardev pipe,id=stdin,path="$tmpdir/stdin" -device virtserialport,chardev=stdin,nr=1 \
  -chardev pipe,id=stdout,path="$tmpdir/stdout" -device virtserialport,chardev=stdout,nr=2 \
  -chardev pipe,id=stderr,path="$tmpdir/stderr" -device virtserialport,chardev=stderr,nr=3 \
  -chardev pipe,id=control,path="$tmpdir/control" -device virtserialport,chardev=control,nr=4 \
  >/dev/null 2>&1 &

send_payload "$@" > "$tmpdir/control.in" &

read exit_code < "$tmpdir/control.out"
exit "$exit_code"
