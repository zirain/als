
IMAGE_PREFIX ?= zirain/
APP_NAME ?= als
TAG ?= latest

BUILDX_PLATFORMS := linux/amd64,linux/arm64

.PHONY: docker-buildx
docker-buildx:
	docker buildx build . -t $(IMAGE_PREFIX)$(APP_NAME):$(TAG) --build-arg GO_LDFLAGS="$(GO_LDFLAGS)" --load

.PHONY: docker-push
docker-push:
	docker buildx build . -t $(IMAGE_PREFIX)$(APP_NAME):$(TAG) --build-arg GO_LDFLAGS="$(GO_LDFLAGS)" --push --platform $(BUILDX_PLATFORMS)
