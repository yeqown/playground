.PHONY: build-image

VERSION ?= 0.0.1
IMAGE_NAME ?= docker.io/yeqown/istio-envoy-filter-demo:$(VERSION)

build-image:
	nerdctl build -t $(IMAGE_NAME) . --build-arg VERSION=$(VERSION)

push-image:
	nerdctl push $(IMAGE_NAME)

run:
	nerdctl run -it --rm -p 8080:8080 $(IMAGE_NAME)