PROJECT_NAME := Pulumi Provider Boilerplate

PACK             := provider-boilerplate
PACKDIR          := sdk
PROJECT          := github.com/pulumi/pulumi-provider-boilerplate
NODE_MODULE_NAME := @pulumi/boilerplate
NUGET_PKG_NAME   := Pulumi.Boilerplate

PROVIDER        := pulumi-resource-${PACK}
PROVIDER_PATH   := provider
VERSION_PATH    := ${PROVIDER_PATH}/version.Version

PULUMI          := .pulumi/bin/pulumi

SCHEMA_FILE     := provider/cmd/pulumi-resource-provider-boilerplate/schema.json
export GOPATH   := $(shell go env GOPATH)

WORKING_DIR     := $(shell pwd)
TESTPARALLELISM := 4

prepare:
	@if test -z "${NAME}"; then echo "NAME not set"; exit 1; fi
	@if test -z "${REPOSITORY}"; then echo "REPOSITORY not set"; exit 1; fi
	@if test -z "${ORG}"; then echo "ORG not set"; exit 1; fi
	@if test ! -d "provider/cmd/pulumi-resource-provider-boilerplate"; then "Project already prepared"; exit 1; fi # SED_SKIP

	# SED needs to not fail when encountering unicode characters
	LC_CTYPE=C 
	LANG=C

	mv "provider/cmd/pulumi-resource-provider-boilerplate" provider/cmd/pulumi-resource-${NAME} # SED_SKIP
	
	# In MacOS the -i parameter needs an empty  to execute in place.
	if [[ "${OS}" == "Darwin" ]]; then \
		find . \( -path './.git' -o -path './sdk' \) -prune -o -not -name 'go.sum' -type f -exec sed -i '' '/SED_SKIP/!s,github.com/pulumi/pulumi-[x]yz,${REPOSITORY},g' {} \; ; \
		find . \( -path './.git' -o -path './sdk' \) -prune -o -not -name 'go.sum' -type f -exec sed -i '' '/SED_SKIP/!s/[xX]yz/${NAME}/g' {} \; ; \
		find . \( -path './.git' -o -path './sdk' \) -prune -o -not -name 'go.sum' -type f -exec sed -i '' '/SED_SKIP/!s/[aA]bc/${ORG}/g' {} \; ; \
	else \
		find . \( -path './.git' -o -path './sdk' \) -prune -o -not -name 'go.sum' -type f -exec sed -i '/SED_SKIP/!s,github.com/pulumi/pulumi-[x]yz,${REPOSITORY},g' {} \; ; \
		find . \( -path './.git' -o -path './sdk' \) -prune -o -not -name 'go.sum' -type f -exec sed -i '/SED_SKIP/!s/[xX]yz/${NAME}/g' {} \; ; \
		find . \( -path './.git' -o -path './sdk' \) -prune -o -not -name 'go.sum' -type f -exec sed -i '/SED_SKIP/!s/[aA]bc/${ORG}/g' {} \; ; \
	fi

# Override during CI using `make [TARGET] PROVIDER_VERSION=""` or by setting a PROVIDER_VERSION environment variable
# Local & branch builds will just used this fixed default version unless specified
PROVIDER_VERSION ?= 1.0.0-alpha.0+dev
# Use this normalised version everywhere rather than the raw input to ensure consistency.
VERSION_GENERIC = $(shell pulumictl convert-version --language generic --version "$(PROVIDER_VERSION)")

# Need to pick up locally pinned pulumi-langage-* plugins.
export PULUMI_IGNORE_AMBIENT_PLUGINS = true

ensure::
	go mod tidy

$(SCHEMA_FILE): provider $(PULUMI)
	$(PULUMI) package get-schema $(WORKING_DIR)/bin/${PROVIDER} | \
		jq 'del(.version)' > $(SCHEMA_FILE)

# Codegen generates the schema file and *generates* all sdks. This is a local process and
# does not require the ability to build all SDKs.
#
# To build the SDKs, use `make build_sdks`
#
# Required by CI (weekly-pulumi-update)
codegen: $(SCHEMA_FILE) sdk/dotnet sdk/go sdk/nodejs sdk/python sdk/java

.PHONY: sdk/%
sdk/%: $(SCHEMA_FILE)
	rm -rf $@
	$(PULUMI) package gen-sdk --language $* $(SCHEMA_FILE) --version "${VERSION_GENERIC}"

sdk/java: $(SCHEMA_FILE)
	rm -rf $@
	$(PULUMI) package gen-sdk --language java $(SCHEMA_FILE)

sdk/python: $(SCHEMA_FILE)
	rm -rf $@
	$(PULUMI) package gen-sdk --language python $(SCHEMA_FILE) --version "${VERSION_GENERIC}"
	cp README.md ${PACKDIR}/python/

sdk/dotnet: $(SCHEMA_FILE)
	rm -rf $@
	$(PULUMI) package gen-sdk --language dotnet $(SCHEMA_FILE) --version "${VERSION_GENERIC}"


sdk/go: ${SCHEMA_FILE}
	rm -rf $@
	$(PULUMI) package gen-sdk --language go ${SCHEMA_FILE} --version "${VERSION_GENERIC}"
	cp go.mod ${PACKDIR}/go/pulumi-${PACK}/go.mod
	cd ${PACKDIR}/go/pulumi-${PACK} && \
		go mod edit -module=github.com/pulumi/pulumi-${PACK}/${PACKDIR}/go/pulumi-${PACK} && \
		go mod tidy

.PHONY: provider
provider: bin/${PROVIDER} bin/pulumi-gen-${PACK} # Required by CI

bin/${PROVIDER}:
	cd provider && go build -o $(WORKING_DIR)/bin/${PROVIDER} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION_GENERIC}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER)

.PHONY: provider_debug
provider_debug:
	(cd provider && go build -o $(WORKING_DIR)/bin/${PROVIDER} -gcflags="all=-N -l" -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION_GENERIC}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

test_provider:
	cd provider && go test -short -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM} -coverprofile="coverage.txt" ./...

dotnet_sdk: sdk/dotnet
	cd ${PACKDIR}/dotnet/&& \
		echo "${VERSION_GENERIC}" > version.txt && \
		dotnet build

go_sdk:	sdk/go

nodejs_sdk: sdk/nodejs
	cd ${PACKDIR}/nodejs/ && \
		yarn install && \
		yarn run tsc
	cp README.md LICENSE ${PACKDIR}/nodejs/package.json ${PACKDIR}/nodejs/yarn.lock ${PACKDIR}/nodejs/bin/

python_sdk: sdk/python
	cp README.md ${PACKDIR}/python/
	cd ${PACKDIR}/python/ && \
		rm -rf ./bin/ ../python.bin/ && cp -R . ../python.bin && mv ../python.bin ./bin && \
		python3 -m venv venv && \
		./venv/bin/python -m pip install build && \
		cd ./bin && \
		../venv/bin/python -m build .

java_sdk:: PACKAGE_VERSION := $(VERSION_GENERIC)
java_sdk:: sdk/java
	cd sdk/java/ && \
		gradle --console=plain build

.PHONY: build
build:: provider build_sdks

.PHONY: build_sdks
build_sdks: dotnet_sdk go_sdk nodejs_sdk python_sdk java_sdk

# Required for the codegen action that runs in pulumi/pulumi
only_build:: build

lint:
	golangci-lint --path-prefix provider --config .golangci.yml run --fix


install:: install_nodejs_sdk install_dotnet_sdk
	cp $(WORKING_DIR)/bin/${PROVIDER} ${GOPATH}/bin


GO_TEST := go test -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM}

test_all:: test
	cd provider/pkg && $(GO_TEST) ./...
	cd tests/sdk/nodejs && $(GO_TEST) ./...
	cd tests/sdk/python && $(GO_TEST) ./...
	cd tests/sdk/dotnet && $(GO_TEST) ./...
	cd tests/sdk/go && $(GO_TEST) ./...

install_dotnet_sdk::
	rm -rf $(WORKING_DIR)/nuget/$(NUGET_PKG_NAME).*.nupkg
	mkdir -p $(WORKING_DIR)/nuget
	find . -name '*.nupkg' -print -exec cp -p {} ${WORKING_DIR}/nuget \;

install_python_sdk::
	#target intentionally blank

install_go_sdk::
	#target intentionally blank

install_java_sdk::
	#target intentionally blank

install_nodejs_sdk::
	-yarn unlink --cwd $(WORKING_DIR)/sdk/nodejs/bin
	yarn link --cwd $(WORKING_DIR)/sdk/nodejs/bin

test:: test_provider
	cd examples && go test -v -tags=all -timeout 2h

$(PULUMI): HOME := $(WORKING_DIR)
$(PULUMI): go.mod
	@ PULUMI_VERSION="$$(cd provider && go list -m github.com/pulumi/pulumi/pkg/v3 | awk '{print $$2}')"; \
	if [ -x $(PULUMI) ]; then \
		CURRENT_VERSION="$$($(PULUMI) version)"; \
		if [ "$${CURRENT_VERSION}" != "$${PULUMI_VERSION}" ]; then \
			echo "Upgrading $(PULUMI) from $${CURRENT_VERSION} to $${PULUMI_VERSION}"; \
			rm $(PULUMI); \
		fi; \
	fi; \
	if ! [ -x $(PULUMI) ]; then \
		curl -fsSL https://get.pulumi.com | sh -s -- --version "$${PULUMI_VERSION#v}"; \
	fi

# Set these variables to enable signing of the windows binary
AZURE_SIGNING_CLIENT_ID ?=
AZURE_SIGNING_CLIENT_SECRET ?=
AZURE_SIGNING_TENANT_ID ?=
AZURE_SIGNING_KEY_VAULT_URI ?=
SKIP_SIGNING ?=

bin/jsign-6.0.jar:
	mkdir -p bin
	wget https://github.com/ebourg/jsign/releases/download/6.0/jsign-6.0.jar --output-document=bin/jsign-6.0.jar

sign-goreleaser-exe-amd64: GORELEASER_ARCH := amd64_v1
sign-goreleaser-exe-arm64: GORELEASER_ARCH := arm64

# Set the shell to bash to allow for the use of bash syntax.
sign-goreleaser-exe-%: SHELL:=/bin/bash
sign-goreleaser-exe-%: bin/jsign-6.0.jar
	@# Only sign windows binary if fully configured.
	@# Test variables set by joining with | between and looking for || showing at least one variable is empty.
	@# Move the binary to a temporary location and sign it there to avoid the target being up-to-date if signing fails.
	@set -e; \
	if [[ "${SKIP_SIGNING}" != "true" ]]; then \
		if [[ "|${AZURE_SIGNING_CLIENT_ID}|${AZURE_SIGNING_CLIENT_SECRET}|${AZURE_SIGNING_TENANT_ID}|${AZURE_SIGNING_KEY_VAULT_URI}|" == *"||"* ]]; then \
			echo "Can't sign windows binaries as required configuration not set: AZURE_SIGNING_CLIENT_ID, AZURE_SIGNING_CLIENT_SECRET, AZURE_SIGNING_TENANT_ID, AZURE_SIGNING_KEY_VAULT_URI"; \
			echo "To rebuild with signing delete the unsigned windows exe file and rebuild with the fixed configuration"; \
			if [[ "${CI}" == "true" ]]; then exit 1; fi; \
		else \
			file=dist/build-provider-sign-windows_windows_${GORELEASER_ARCH}/pulumi-resource-provider-boilerplate.exe; \
			mv $${file} $${file}.unsigned; \
			az login --service-principal \
				--username "${AZURE_SIGNING_CLIENT_ID}" \
				--password "${AZURE_SIGNING_CLIENT_SECRET}" \
				--tenant "${AZURE_SIGNING_TENANT_ID}" \
				--output none; \
			ACCESS_TOKEN=$$(az account get-access-token --resource "https://vault.azure.net" | jq -r .accessToken); \
			java -jar bin/jsign-6.0.jar \
				--storetype AZUREKEYVAULT \
				--keystore "PulumiCodeSigning" \
				--url "${AZURE_SIGNING_KEY_VAULT_URI}" \
				--storepass "$${ACCESS_TOKEN}" \
				$${file}.unsigned; \
			mv $${file}.unsigned $${file}; \
			az logout; \
		fi; \
	fi

# To make an immediately observable change to .ci-mgmt.yaml:
#
# - Edit .ci-mgmt.yaml
# - Run make ci-mgmt to apply the change locally.
#
ci-mgmt: .ci-mgmt.yaml
	go run github.com/pulumi/ci-mgmt/provider-ci@master generate
.PHONY: ci-mgmt

.PHONY:local_generate
local_generate: # Required by CI

.PHONY: generate_schema
generate_schema: ${SCHEMA_PATH} # Required by CI

.PHONY: build_go install_go_sdk
generate_go: sdk/go # Required by CI
build_go: # Required by CI

.PHONY: build_java install_java_sdk
generate_java: sdk/java # Required by CI
build_java: java_sdk # Required by CI

.PHONY: build_python install_python_sdk
generate_python: sdk/python # Required by CI
build_python: python_sdk # Required by CI

.PHONY: build_nodejs install_nodejs_sdk
generate_nodejs: sdk/nodejs # Required by CI
build_nodejs: nodejs_sdk # Required by CI

.PHONY: build_dotnet install_dotnet_sdk
generate_dotnet: sdk/dotnet # Required by CI
build_dotnet: dotnet_sdk # Required by CI

bin/pulumi-gen-${PACK}: # Required by CI
	touch bin/pulumi-gen-${PACK}
