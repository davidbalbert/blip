# debian:13.1 on 2025-11-10
FROM debian@sha256:01a723bf5bfb21b9dda0c9a33e0538106e4d02cce8f557e118dd61259553d598

ARG DEBIAN_SNAPSHOT=20251110T143905Z
ARG SUITE=trixie

ENV TZ=UTC \
    LANG=C.UTF-8 \
    LC_ALL=C.UTF-8 \
    DEBIAN_FRONTEND=noninteractive

RUN <<EOF
set -eux
cat > /etc/apt/sources.list.d/debian.sources <<SOURCES
Types: deb
URIs: http://snapshot.debian.org/archive/debian/${DEBIAN_SNAPSHOT}
Suites: ${SUITE}
Components: main
Check-Valid-Until: no
Signed-By: /usr/share/keyrings/debian-archive-keyring.gpg

Types: deb
URIs: http://snapshot.debian.org/archive/debian-security/${DEBIAN_SNAPSHOT}
Suites: ${SUITE}-security
Components: main
Check-Valid-Until: no
Signed-By: /usr/share/keyrings/debian-archive-keyring.gpg
SOURCES

apt -o Acquire::Check-Valid-Until=false update
apt install -y bc bison build-essential curl flex xz-utils
rm -rf /var/lib/apt/lists/*
EOF

WORKDIR /build

ENTRYPOINT ["/build/build.sh"]
