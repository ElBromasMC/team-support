# Include .env file if exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Variable for the environment mode 'production' or 'development'
ENV ?= production

TMP_DIR := ./tmp
BUILD_DIR := ./build
CMD_DIR := ./cmd
VIEW_DIR := ./view
JS_DIR := $(VIEW_DIR)/js

GO_FILES := $(shell find . -path ./node_modules -prune -o -path $(VIEW_DIR) -prune -o -path $(CMD_DIR) -prune -o -name '*.go' -print)
TEMPL_FILES := $(shell find $(VIEW_DIR) -name '*.templ')
TEMPL_GO_FILES = $(TEMPL_FILES:.templ=_templ.go)
JS_FILES := $(shell find $(JS_DIR) -name '*.js')

MAINJS_OUT := ./assets/static/js/main.js
TAILWINDCSS_OUT := ./assets/static/css/tailwind.css

TAILWINDCSS_LOG := tailwind.log

$(TMP_DIR)/%: $(BUILD_DIR)/% $(MAINJS_OUT) $(TAILWINDCSS_OUT)
	touch "$@"

# The final build step
ifeq ($(ENV),development)
$(BUILD_DIR)/%: $(wildcard $(CMD_DIR)/%/*) $(GO_FILES) $(TEMPL_GO_FILES)
	go build -tags=dev -o "$@" "$(CMD_DIR)/$*"
else
$(BUILD_DIR)/%: $(wildcard $(CMD_DIR)/%/*) $(GO_FILES) $(TEMPL_GO_FILES) $(MAINJS_OUT) $(TAILWINDCSS_OUT)
	go build -o "$@" "$(CMD_DIR)/$*"
endif

# Build step for templ source
%_templ.go: %.templ
	templ generate -f "$<"

# Build step for main.js
$(MAINJS_OUT): $(JS_FILES)
ifeq ($(ENV),development)
	npx esbuild "$(JS_DIR)/main.js" --outfile="$@" --bundle
else
	npx esbuild "$(JS_DIR)/main.js" --outfile="$@" --bundle --minify
endif

# Build step for tailwind.css
$(TAILWINDCSS_OUT): $(TEMPL_FILES) tailwind.config.js tailwind.css
ifeq ($(ENV),development)
	./dev/tailwind.sh "$(TMP_DIR)" "$(TAILWINDCSS_LOG)"
else
	npx tailwindcss build -i tailwind.css -o "$@" --minify
endif
	touch "$@"

# Live reload
.PHONY: live
live:
ifeq ($(ENV),development)
	mkdir -p "$(TMP_DIR)"
	npx tailwindcss build -i tailwind.css -o "$(TAILWINDCSS_OUT)" --watch=always &> "$(TMP_DIR)/$(TAILWINDCSS_LOG)" &
	node ./dev/live-reload.js &
	air
else
	@echo "Must run with ENV=development"
endif

.PHONY: clean
clean:
	@if [ -z "$(VIEW_DIR)" ] || [ -z "$(BUILD_DIR)" ]; then \
		echo "Error: VIEW_DIR or BUILD_DIR is not set."; \
		exit 1; \
	fi
	rm -f "$(VIEW_DIR)"/**/*_templ.go
	rm -f "$(TAILWINDCSS_OUT)" "$(MAINJS_OUT)"
	rm -rf "$(BUILD_DIR)"
	rm -rf "$(TMP_DIR)"

.SECONDARY: