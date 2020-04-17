BUILD_DIR=build
SERVICES=serviceA serviceB serviceC serviceD
DOCKERS=$(addprefix docker_,$(SERVICES))
CGO_ENABLED?=0
GOOS?=linux
REPO=

define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -ldflags "-s -w" -o ${BUILD_DIR}/$(1) cmd/$(1)/main.go
endef

define make_docker
	docker build --build-arg SVC_NAME=$(subst docker_,,$(1)) --tag=$(REPO)/$(subst docker_,,$(1)) -f ./docker/Dockerfile ./build/
endef

all: $(SERVICES)

.PHONY: all $(SERVICES) dockers latest release

clean:
	rm -rf ${BUILD_DIR}

install:
	cp ${BUILD_DIR}/* $(GOBIN)

$(SERVICES):
	$(call compile_service,$(@))

$(DOCKERS):
	$(call make_docker,$(@))

dockers: $(SERVICES) $(DOCKERS)