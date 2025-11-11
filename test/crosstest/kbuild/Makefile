LINUX_VERSION := 6.17.7
LINUX_SHA256 := ddf2ea0d4439e1d57136be3623102af9458f601f5b1cb77e83246e88aea09d0e
LINUX_DATE := 2025-11-02
LINUX_TAR := pkg/linux-$(LINUX_VERSION).tar.xz
LINUX_URL := https://cdn.kernel.org/pub/linux/kernel/v6.x/linux-$(LINUX_VERSION).tar.xz

LINUX_IMAGE_SHA256 := 94d573898fa4f505a2f7fc720fd389bd6e9b5556cf8c3f9dec68893b5ddba96c

BUSYBOX_VERSION := 1.36.1
BUSYBOX_SHA256 := b8cc24c9574d809e7279c3be349795c5d5ceb6fdf19ca709f80cde50e47de314
BUSYBOX_DATE := 2023-05-19
BUSYBOX_TAR := pkg/busybox-$(BUSYBOX_VERSION).tar.bz2
BUSYBOX_URL := https://busybox.net/downloads/busybox-$(BUSYBOX_VERSION).tar.bz2

DOCKER_IMAGE_NAME := kbuild

DOCKER_RUN_ARGS := \
	-v $(PWD)/pkg:/pkg:ro \
	-v $(PWD)/config:/config:ro \
	-v $(PWD)/build.sh:/build/build.sh:ro \
	-v $(PWD)/build-busybox.sh:/build/build-busybox.sh:ro \
	-v $(PWD)/build-linux.sh:/build/build-linux.sh:ro \
	-v $(PWD)/init.sh:/src/init.sh:ro \
	-v $(PWD)/out:/out:rw \
	-e BUSYBOX_VERSION=$(BUSYBOX_VERSION) \
	-e BUSYBOX_DATE=$(BUSYBOX_DATE) \
	-e LINUX_VERSION=$(LINUX_VERSION) \
	-e LINUX_DATE=$(LINUX_DATE)

.DEFAULT_GOAL := out/Image
.PHONY: all shell clean

.container: Dockerfile
	docker build -t $(DOCKER_IMAGE_NAME) .
	@touch .container

$(LINUX_TAR):
	@echo "Downloading linux kernel..."
	curl -L -o $(LINUX_TAR) $(LINUX_URL)

$(BUSYBOX_TAR):
	@echo "Downloading busybox..."
	curl -L -o $(BUSYBOX_TAR) $(BUSYBOX_URL)

.checksums: $(LINUX_TAR) $(BUSYBOX_TAR)
	@echo "Verifying SHA-256 checksums..."
	@cd pkg && echo "$(LINUX_SHA256) linux-$(LINUX_VERSION).tar.xz" | sha256sum -c -
	@cd pkg && echo "$(BUSYBOX_SHA256) busybox-$(BUSYBOX_VERSION).tar.bz2" | sha256sum -c -
	@touch .checksums

out/Image: .container .checksums
	docker run --rm $(DOCKER_RUN_ARGS) $(DOCKER_IMAGE_NAME)
	@echo "Verifying out/Image"
	@cd out && echo "$(LINUX_IMAGE_SHA256) Image" | sha256sum -c -

shell: .container .checksums
	docker run -it --rm $(DOCKER_RUN_ARGS) --entrypoint /bin/bash $(DOCKER_IMAGE_NAME)

clean:
	rm -rf out/* pkg/* .container .checksums
