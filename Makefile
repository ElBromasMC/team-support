# Include .env file if exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

#ES_FLAGS := --minify
#TAILWIND_FLAGS := --minify

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

# The final build step
$(BUILD_DIR)/%: $(CMD_DIR)/%.go $(GO_FILES) $(TEMPL_GO_FILES) $(MAINJS_OUT) $(TAILWINDCSS_OUT)
	go build -o $@ $<

# Build step for templ source
%_templ.go: %.templ
	templ generate -f $<

# Build step for main.js
$(MAINJS_OUT): $(JS_FILES)
	npx esbuild $(JS_DIR)/main.js --outfile=$@ --bundle $(ES_FLAGS)

# Build step for tailwind.css
$(TAILWINDCSS_OUT): $(TEMPL_FILES)
	npx tailwindcss build -i tailwind.css -o $@ $(TAILWIND_FLAGS) && touch $@

# Live reload
.PHONY: live
live:
	node live-reload.js &
	air

.PHONY: clean
clean:
	rm $(VIEW_DIR)/**/*_templ.go
	rm $(TAILWINDCSS_OUT) $(MAINJS_OUT)
	rm $(BUILD_DIR)/*
