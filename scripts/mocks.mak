# Define the rule to run `go generate` in each package
.PHONY: mock clean_mock

# Get a list of packages with changes in Go files
PACKAGES := $(shell go list github.com/cldcvr/terrarium/... | grep -v /vendor/ | grep -v '/mocks$$')
PACKAGE_DIRS := $(PACKAGES:github.com/cldcvr/terrarium/%=./%)

mock: $(addsuffix /mocks, $(PACKAGE_DIRS))  ## generate mock files for updated go packages

clean_mock:  ## Delete all the generated mock files
	@echo -e "Removing Mocks"
	@find . -type d -name mocks -prune -exec rm -rf {} \;

%/mocks: %/*.go
	@echo "generating mocks for $*"
	@(cd $* && rm -rf mocks && go generate && mkdir -p mocks)
