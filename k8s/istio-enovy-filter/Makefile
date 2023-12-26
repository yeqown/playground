.PHONY: build-image push-image run

VERSION ?= 0.0.1
IMAGE_NAME ?= docker.io/yeqown/istio-envoy-filter-demo:$(VERSION)

build-image:
	nerdctl.lima build -t $(IMAGE_NAME) . --build-arg VERSION=$(VERSION)

push-image:
	nerdctl.lima push $(IMAGE_NAME)

run:
	nerdctl.lima run -it --rm -p 8080:8080 $(IMAGE_NAME)