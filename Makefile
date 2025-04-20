# Exporting bin folder to the path for makefile
export PATH   := $(PWD)/bin:$(PATH)
# Default Shell
export SHELL  := bash
# Type of OS: Linux or Darwin.
export OSTYPE := $(shell uname -s | tr A-Z a-z)
export ARCH := $(shell uname -m)

# --- Tooling & Variables ----------------------------------------------------------------
include ./scripts/make/tools.Makefile
include ./scripts/make/help.Makefile

install-deps: migrate air gotestsum tparse mockery ## Install Development Dependencies (localy).
deps: $(MIGRATE) $(AIR) $(GOTESTSUM) $(TPARSE) $(MOCKERY) $(GOLANGCI) ## Checks for Global Development Dependencies.
deps:
	@echo "Required Tools Are Available"

dev-air: $(AIR) ## Starts AIR (Continuous Development app).
ifeq ($(OSTYPE), linux)
	air -c .air.linux.toml
else ifeq ($(OSTYPE), windows)
	air -c .air.windows.toml
endif

lint: $(GOLANGCI) ## Runs golangci-lint with predefined configuration
	@echo "Applying linter"
	golangci-lint version
	golangci-lint run -c .golangci.yaml ./...


# ~~~ Cleans ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

clean: clean-artifacts clean-docker clean-deps ## Execute all clean commands below

clean-artifacts: ## Removes Artifacts (*.out)
	@printf "Cleanning artifacts... "
	@rm -f *.out
	@echo "done."

clean-docker: ## Removes dangling docker images
	@ docker image prune -f

clean-deps: ## Cleans up the local dependencies
	@printf "Cleaning dependencies... "
	@rm -rf bin
	@echo "done."