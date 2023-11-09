REQ_ENVS := dev staging prod
REQUIREMENTS_DIR := ./requirements
PLATFORM_DIR := ./platforms
TARGET_DIR := ./layers

PLATFORM1_DIR := $(PLATFORM_DIR)/01-common
PLATFORM2_DIR := $(PLATFORM_DIR)/02-team-shared
PLATFORM3_DIR := $(PLATFORM_DIR)/03-team-env

TARGET_COMM_DIR := $(TARGET_DIR)/01-common
TARGET_BU_DIR := $(TARGET_DIR)/02-bu

REQ_ENV_YAML := $(REQUIREMENTS_DIR)/.env.yaml # this file is auto generated based on number of envs given on top here
REQ_BASELINE_YAML := $(REQUIREMENTS_DIR)/common.yaml
REQ_TEAMS_YAML := $(shell find $(REQUIREMENTS_DIR) -name 'team-*.yaml')
REQ_SHARED_TEAMS_YAML := $(shell find $(REQUIREMENTS_DIR) -name 'shared-team-*.yaml')
REQ_ENV_TEAMS_YAML := $(shell find $(REQUIREMENTS_DIR) -name 'env-team-*.yaml')
REQ_ALL_YAML := $(REQ_BASELINE_YAML) $(REQ_ENV_YAML) $(REQ_SHARED_TEAMS_YAML) $(REQ_ENV_TEAMS_YAML) $(REQ_TEAMS_YAML)
REQ_TEAM_NAMES := $(REQ_TEAMS_YAML:$(REQUIREMENTS_DIR)/team-%.yaml=%)

TARGET_TEAM_SHARED_DIR := $(REQ_TEAM_NAMES:%=$(TARGET_BU_DIR)/%/shared)
TARGET_TEAM_ENV_DIR := $(foreach env,$(REQ_ENVS),$(REQ_TEAM_NAMES:%=$(TARGET_BU_DIR)/%/$(env)))
TARGETS_ALL := $(TARGET_COMM_DIR) $(TARGET_TEAM_SHARED_DIR) $(TARGET_TEAM_ENV_DIR)

PLATFORMS_ALL := $(PLATFORM1_DIR) $(PLATFORM2_DIR) $(PLATFORM3_DIR)
PLATFORM_YAML_FILES :=  $(PLATFORMS_ALL:%=%/terrarium.yaml)

append_prefix_to_arr = $(1:%=$(2)%) # function to append a prefix to an array. $(1)=array;$(2)=prefix
req_team_env_yaml = $(filter %-$(2).yaml,$(REQ_TEAMS_YAML)) $(filter %-$(1)-$(2).yaml,$(REQ_ENV_TEAMS_YAML)) # function to get yaml files for a specific team env. $(1)=env; $(2)=team name
req_team_shared_yaml = $(filter %-$(1).yaml,$(REQ_TEAMS_YAML) $(REQ_SHARED_TEAMS_YAML)) $(REQ_ENV_YAML) # function to get yaml files for the shared dir in a team $(1)=team name

.PHONY: dryrun
dryrun:
	@echo "Platforms: $(PLATFORMS_ALL:%=\n\t%)\n"
	@echo "Targets: $(TARGETS_ALL:%=\n\t%)\n"
	@echo "Requirements: $(REQ_ALL_YAML:%=\n\t%)\n"

.PHONY: lz_cleanup
lz_cleanup:
	rm -rf $(TARGET_COMM_DIR)
	rm -rf $(TARGET_BU_DIR)
	rm -rf $(REQ_ENV_YAML)

.PHONY: lz
lz: $(TARGETS_ALL)

.PHONY: lint_all
lint_all: $(PLATFORM_YAML_FILES)

.PHONY: lint_cleanup
lint_cleanup:
	rm  -f $(PLATFORM_YAML_FILES)

# generate the manifest file for envs based on the set of given env types in `REQ_ENVS`
$(REQ_ENV_YAML):
	@echo "id: env_deps" > $@
	@echo "dependencies:" >> $@
	@$(foreach env,$(REQ_ENVS),echo "- id: $(env)" >> $@; echo "  use: environment" >> $@;)

# generate terraform code for the baseline 01 layer
$(TARGET_COMM_DIR): $(REQ_ALL_YAML) $(PLATFORM1_DIR)/terrarium.yaml
	rm -rf $@
	mkdir -p $@
	terrarium generate --skip-env-file -p $(PLATFORM1_DIR) -o $@ $(call append_prefix_to_arr,$(REQ_ALL_YAML),-a )

# generate terraform code for a specific team's shared dir
$(TARGET_BU_DIR)/%/shared: $(call req_team_shared_yaml,%) $(PLATFORM2_DIR)/terrarium.yaml
	rm -rf $@
	mkdir -p $@
	terrarium generate --skip-env-file -p $(PLATFORM2_DIR) -o $@ $(call append_prefix_to_arr,$(call req_team_shared_yaml,$*),-a )

# generate terraform code for a specific team and specific env
define rules_team_env
$(TARGET_BU_DIR)/%/$(1): $(call req_team_env_yaml,$(1),%) $(PLATFORM3_DIR)/terrarium.yaml
	rm -rf $$@
	mkdir -p $$@
	terrarium generate --skip-env-file -p $$(PLATFORM3_DIR) -o $$@ -c $(1) $$(call append_prefix_to_arr,$$(call req_team_env_yaml,$(1),$$*),-a )
endef

$(foreach env,$(REQ_ENVS),$(eval $(call rules_team_env,$(env))))

$(PLATFORM_DIR)/%/terrarium.yaml: $(PLATFORM_DIR)/%/*.tf
	rm -f $@
	terrarium platform lint -d $(dir $@)
