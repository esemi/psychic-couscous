# Package dir
GO_DIR ?= $(CURDIR)

# Package import path
GO_PKG ?= $(shell go list -e -m -f "{{ .Dir }}" 2> /dev/null || go list -e -f "{{ .ImportPath }}" 2> /dev/null || pwd)

# Package version
GO_VER ?= $(shell date -u +%Y-%m-%d.%H:%M:%S)

# Package binary prefix
GO_BIN ?= $(shell basename $$(dirname $(GO_PKG)))

# Package binnary delimiter
GO_DEL ?= -

# Set the mode for code-coverage
GO_TEST_COVERAGE_MODE ?= count
GO_TEST_COVERAGE_FILE_NAME ?= coverage.out

# Set a default `min_confidence` value for `golint`
GO_LINT_MIN_CONFIDENCE ?= 0.2

# .env
-include Makefile.env

all: help

.PHONY: help # Show list of targets with description. You're looking at it
help:
	@grep "^.PHONY: .* #" Makefile | sed "s/\.PHONY: \(.*\) # \(.*\)/\1 - \2/" | expand -t20

.PHONY: build # Compile application source code to binaries files
build:
	@echo "Build package: $(GO_PKG)"
	@ERR=0; \
	for CMD in $$(find "./cmd" -maxdepth 1 -mindepth 1 -type d -print); do \
		BIN=$$(basename "$${CMD}"); \
		go build -gcflags="-trimpath=$(GO_DIR)" \
				 -asmflags="-trimpath=$(GO_DIR)" \
				 -ldflags "-X main.Version=$(GO_VER)" \
				 -o "$(GO_DIR)/.bin/$${BIN}/$(GO_BIN)$(GO_DEL)$${BIN}" "$${CMD}" || { \
			ERR=$$?; \
			break; \
		}; \
	done; \
	if [ $$ERR != 0 ]; then \
		exit $$ERR; \
	fi

.PHONY: generate # Generate auto generated code
generate:
	@echo "Run easyjson"
	@go list -f "{{ .Dir }}" ./... | xargs -I "{}" grep -lrw "{}" -e "easyjson:json" | sort | uniq | xargs -I "{}" easyjson "{}" || exit 1

.PHONY: compress # Compress application binaries
compress:
	@echo "Compress binaries"
	@ERR=0; \
    for CMD in $$(find "./cmd" -maxdepth 1 -mindepth 1 -type d -print); do \
    	BIN=$$(basename "$${CMD}"); \
    	upx --best -q "$(GO_DIR)/.bin/$${BIN}/$(GO_BIN)$(GO_DEL)$${BIN}" || { \
    		ERR=$$?; \
    		break; \
    	}; \
    done; \
    if [ $$ERR != 0 ]; then \
    	exit $$ERR; \
    fi;

.PHONY: conf-gomod # Setup GOPROVATE and configure GIT
conf-gomod:
	@go env -w GOPRIVATE=*.mycompany.com
	@git config --global url.ssh://git@git.mycompany.com/.insteadOf https://git.mycompany.com/

.PHONY: install # Install dependencies
install: conf-gomod
	@echo "Install dependencies"
	go mod download

.PHONY: update # Update dependencies
update: git
	@echo "Update dependencies"
	go get -u -t all
	@echo "Optimize dependencies"
	go mod tidy -v

.PHONY: test # Run application tests
test:
	@echo "Run unit tests"
	@go test -v ./...

.PHONY: test-with-coverage # Run application tests with code coverage
test-with-coverage:
	@echo "Run unit tests with coverage"
	@go test -cover ./...

.PHONY: test-with-coverage-profile # Run application tests with code coverage profile
test-with-coverage-profile:
	@echo "Run unit tests with coverage profile"
	@echo "mode: ${GO_TEST_COVERAGE_MODE}" > "${GO_TEST_COVERAGE_FILE_NAME}"
	@go test -coverpkg=`go list ./... | grep -vE 'command|domain' | tr '\n' ','` -covermode ${GO_TEST_COVERAGE_MODE} -coverprofile=${GO_TEST_COVERAGE_FILE_NAME} ./...
	@echo "Generate coverage report";
	@go tool cover -func="${GO_TEST_COVERAGE_FILE_NAME}";
	@rm "${GO_TEST_COVERAGE_FILE_NAME}";

.PHONY: fix # Fix code style
fix: fix-format fix-import

.PHONY: fix-format # Fix formatting of code
fix-format:
	@echo "Fix formatting"
	@gofmt -w ${GO_FMT_FLAGS} $$(go list -f "{{ .Dir }}" ./...); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

.PHONY: fix-import # Fix code style of imports
fix-import:
	@echo "Fix imports"
	@errors=$$(goimports -l -w -local $(GO_PKG) $$(go list -f "{{ .Dir }}" ./...)); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

.PHONY: lint # Lint code style of source code
lint: lint-format lint-import lint-style

.PHONY: lint-format # Lint formatting of source code
lint-format:
	@echo "Check formatting"
	golangci-lint run --disable-all --enable gofmt

.PHONY: lint-import # Lint code style of imports
lint-import:
	@echo "Check imports"
	golangci-lint run --disable-all --enable goimports

.PHONY: lint-style # Lint code style of source code
lint-style:
	@echo "Check code style"
	@errors=$$(golint -min_confidence=${GO_LINT_MIN_CONFIDENCE} $$(go list ./...)); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

.PHONY: clean # Delete auto generated files
clean:
	@echo "Cleanup"
	@find . -type f -name "*easyjson*" -delete
	@find . -type f -name "*coverage*.out" -delete
