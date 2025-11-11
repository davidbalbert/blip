LINUX_VERSION := 6.17.7
LINUX_SHA256 := ddf2ea0d4439e1d57136be3623102af9458f601f5b1cb77e83246e88aea09d0e
LINUX_DATE := 2025-11-02
LINUX_TAR := pkg/linux-$(LINUX_VERSION).tar.xz
LINUX_URL := https://cdn.kernel.org/pub/linux/kernel/v6.x/linux-$(LINUX_VERSION).tar.xz

BUSYBOX_VERSION := 1.36.1
BUSYBOX_SHA256 := b8cc24c9574d809e7279c3be349795c5d5ceb6fdf19ca709f80cde50e47de314
BUSYBOX_DATE := 2023-05-19
BUSYBOX_TAR := pkg/busybox-$(BUSYBOX_VERSION).tar.bz2
BUSYBOX_URL := https://busybox.net/downloads/busybox-$(BUSYBOX_VERSION).tar.bz2

IMAGE_NAME := kbuild

.PHONY: all build verify shell clean

all: build verify

build:
	docker build -t $(IMAGE_NAME) .

$(LINUX_TAR):
	@echo "Downloading linux kernel..."
	curl -L -o $(LINUX_TAR) $(LINUX_URL)

$(BUSYBOX_TAR):
	@echo "Downloading busybox..."
	curl -L -o $(BUSYBOX_TAR) $(BUSYBOX_URL)

verify: $(LINUX_TAR) $(BUSYBOX_TAR)
	@echo "Verifying SHA-256 checksums..."
	@cd pkg && echo "$(LINUX_SHA256) linux-$(LINUX_VERSION).tar.xz" | sha256sum -c -
	@cd pkg && echo "$(BUSYBOX_SHA256) busybox-$(BUSYBOX_VERSION).tar.bz2" | sha256sum -c -

shell: verify build
	docker run -it --rm \
		-v $(PWD)/pkg:/pkg:ro \
		-v $(PWD)/config:/config:ro \
		-v $(PWD)/build-busybox.sh:/build/build-busybox.sh:ro \
		-v $(PWD)/build-linux.sh:/build/build-linux.sh:ro \
		-v $(PWD)/init.sh:/src/init.sh:ro \
		-v $(PWD)/out:/out:rw \
		-e BUSYBOX_VERSION=$(BUSYBOX_VERSION) \
		-e BUSYBOX_DATE=$(BUSYBOX_DATE) \
		-e LINUX_VERSION=$(LINUX_VERSION) \
		-e LINUX_DATE=$(LINUX_DATE) \
		$(IMAGE_NAME)

clean:
	rm -rf src
