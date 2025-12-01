#!/bin/sh
mount -t proc proc /proc
mount -t sysfs sysfs /sys

stdin=/dev/vport0p1
stdout=/dev/vport0p2
stderr=/dev/vport0p3
control=/dev/vport0p4

# Strip leading zeros to avoid parsing as octal
strip_zeros() {
  echo "$1" | sed 's/^0*//' | { read n; echo "${n:-0}"; }
}

exec 3<>$control

# Protocol: [8-byte total_size][8-byte program_size][program][8-byte args_size][args]
total_size=$(strip_zeros "$(dd bs=1 count=8 <&3 2>/dev/null)")
dd bs=4096 count=$(( (total_size + 4095) / 4096 )) <&3 2>/dev/null | head -c "$total_size" > /tmp/payload

# Extract program
program_size=$(strip_zeros "$(head -c 8 /tmp/payload)")
tail -c +9 /tmp/payload | head -c "$program_size" > /tmp/program
chmod +x /tmp/program

# Extract args
args_offset=$((8 + program_size))
args_size=$(strip_zeros "$(tail -c +$((args_offset + 1)) /tmp/payload | head -c 8)")
if [ "$args_size" -gt 0 ]; then
  tail -c +$((args_offset + 9)) /tmp/payload | head -c "$args_size" > /tmp/args
else
  : > /tmp/args
fi

# Build arg list (null separated input)
set --
while IFS= read -r -d '' arg; do
  set -- "$@" "$arg"
done < /tmp/args


/tmp/program "$@" <$stdin >$stdout 2>$stderr
echo $? >&3

echo o > /proc/sysrq-trigger
