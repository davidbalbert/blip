qemu-system-aarch64 \
  -machine virt,gic-version=3 \
  -cpu host \
  -accel hvf \
  -nographic \
  -m 1G \
  -kernel ./out/arm64/Image \
  -append "console=ttyAMA0 loglevel=0"
