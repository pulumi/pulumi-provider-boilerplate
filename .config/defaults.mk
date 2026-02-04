# This file contains default make targets used by provders. To use it,
# add this to the top of your Makefile:
#
#     include .config/defaults.mk
#
# This will expose default targets that should "just work" without any
# customization.
#
# To extend or customize the default behavior, you can:
#
#    1. Override targets:
#
#           lint:
#              ./custom-lint.sh
#
#    2. Compose targets:
#
#           lint: lint.pre default.lint lint.post
#
#       (Remember make processes dependencies from left to right.)
#
# Default targets contain our actual implementations and are prefixed with
# "default." in order to enable wrapping. They are not meant to be overridden.

# TODO(https://github.com/pulumi/ci-mgmt/issues/2033): This should lint the
# entire repo, not just the provider subdirectory.
.PHONY: default.lint
default.lint:
	cd provider && golangci-lint run --path-prefix provider -c ../.golangci.yml

.PHONY: default.lint.fix
default.lint.fix:
	cd provider && golangci-lint run --fix --path-prefix provider -c ../.golangci.yml

# The public targets below are actually consumed by CI.

.PHONY: lint
lint: default.lint

.PHONY: lint.fix
lint.fix: default.lint.fix
