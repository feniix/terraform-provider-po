NAME=terraform-provider-po
VERSION ?= $(shell cat VERSION)

default: install

build: build/linux build/osx

build/linux:
	GOOS=linux GOARCH=amd64 go build -o out/${NAME}-linux-amd64 ./cmd/${NAME}

build/osx:
	GOOS=darwin GOARCH=amd64 go build -o out/${NAME}-darwin-amd64 ./cmd/${NAME}

install: build/osx
	mkdir -p ~/.terraform.d/plugins/github.com/feniix/po/${VERSION}/darwin_amd64
	mv ./out/${NAME}-darwin-amd64 ~/.terraform.d/plugins/github.com/feniix/po/${VERSION}/darwin_amd64/${NAME}

clean:
	rm -rf out

clean-tfstate:
	rm -rf ./examples/**/terraform.* ./exa/mples/**.terraform ../examples/**/.terraform.lock.hcl ./examples/**/crash.log

install-operator:
	kubectl apply -f config/bundle.yaml

setup-cluster:
	kind create cluster --config ./config/kind-config.yaml
	$(MAKE) install-operator


