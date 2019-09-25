# Image URL to use all building/pushing image targets
IMG ?= zhuxiaoyang/ks-scheduler:v1

# Build manager binary
manager: fmt vet
	GOOS=linux GOARCH=amd64 go build -mod=vendor -a -o ks-scheduler /root/cmd

# Build the docker image
docker-build:
	docker build -t $(IMG) .
	docker push $(IMG)

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# Run tests
test: fmt vet
	export KUBEBUILDER_CONTROLPLANE_START_TIMEOUT=1m; ginkgo -v -cover ./pkg/...

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	go run ./cmd/manager/main.go

e2e-test:
	./hack/e2etest.sh

install-travis:
	chmod +x ./hack/*.sh
	./hack/install_tools.sh

.PHONY : clean test


