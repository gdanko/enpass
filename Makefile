GOPATH := $(shell go env GOPATH)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
ENPASS_VERSION := "0.2.3"

GOOS ?= $(shell uname | tr '[:upper:]' '[:lower:]')
GOARCH ?=$(shell arch)

.PHONY: all build install

all: build install

.PHONY: mod-tidy
mod-tidy:
	go mod tidy

.PHONY: build OS ARCH
build: guard-ENPASS_VERSION mod-tidy clean
	@echo "================================================="
	@echo "Building enpass"
	@echo "=================================================\n"

	@if [ ! -d "${GOOS}" ]; then \
		mkdir "${GOOS}"; \
	fi
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o "${GOOS}/enpass"
	sleep 2
	tar -C "${GOOS}" -czvf "enpass_${ENPASS_VERSION}_${GOOS}_${GOARCH}.tgz" enpass; \

.PHONY: clean
clean:
	@echo "================================================="
	@echo "Cleaning enpass"
	@echo "=================================================\n"
	@for OS in darwin linux; do \
		if [ -f $${OS}/enpass ]; then \
			rm -f $${OS}/enpass; \
		fi; \
	done

.PHONY: clean-all
clean-all: clean
	@echo "================================================="
	@echo "Cleaning tarballs"
	@echo "=================================================\n"
	@rm -f *.tgz 2>/dev/null

.PHONY: install
install:
	@echo "================================================="
	@echo "Installing enpass in ${GOPATH}/bin"
	@echo "=================================================\n"

	go install -race

#
# General targets
#
guard-%:
	@if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi
