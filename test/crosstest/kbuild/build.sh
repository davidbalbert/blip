#!/bin/bash
set -euo pipefail

/build/build-busybox.sh
/build/build-linux.sh
