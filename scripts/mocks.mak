# Get a list of packages with changes in Go files
PACKAGES := $(shell go list github.com/cldcvr/terrarium/... | grep -v /vendor/ | grep -v '/mocks$$')
PACKAGE_DIRS := $(PACKAGES:github.com/cldcvr/terrarium/%=./%)


# Define the rule to run `go generate` in each package
.PHONY: mock clean_mock
mock: $(addsuffix /mocks, $(PACKAGE_DIRS))

%/mocks: %/*.go
	@echo "Running go generate in $*..."
	@(cd $* && rm -rf mocks && mkdir -p mocks && go generate)

# Get a list of all 'mocks' directories
MOCKS_DIRS := $(shell find . -type d -name 'mocks')

clean_mock:
	@for dir in $(MOCKS_DIRS); do \
		echo "Removing $$dir..."; \
		rm -rf $$dir; \
	done
