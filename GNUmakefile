VERSION=0.0.1
default: build

PROJECT_NAME := $(notdir $(CURDIR))
GSMTP-CLI := "gsmtp-cli"

.PHONY: $(PROJECT_NAME)

project-name:
	@echo "PROJECT_NAME: $(PROJECT_NAME)"
	@echo "PROJECT_NAME: $(PROJECT_NAME)"

clean: project-name
	rm -rf $(CURDIR)/build

install: clean 
	install -d $(CURDIR)/build/macos_arm64 && install -d $(CURDIR)/build/linux_amd64

build: install
	export GOOS=darwin && export GOARCH=arm64 && go build -C $(CURDIR)/cmd/$(PROJECT_NAME) -o $(CURDIR)/build/macos_arm64
	export GOOS=linux && export GOARCH=amd64 && go build -C $(CURDIR)/cmd/$(PROJECT_NAME) -o $(CURDIR)/build/linux_amd64
	export GOOS=darwin && export GOARCH=arm64 && go build -C $(CURDIR)/cmd/$(GSMTP-CLI) -o $(CURDIR)/build/macos_arm64
	export GOOS=linux && export GOARCH=amd64 && go build -C $(CURDIR)/cmd/$(GSMTP-CLI) -o $(CURDIR)/build/linux_amd64
