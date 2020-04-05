.DEFAULT_GOAL := bin
INSTALL=install -p -m 644

.PHONY: build clean bin release

build: bin

clean:
	@rm -r bin

bin:
	@mkdir -p bin
	@go build -o bin/wejay
	@echo '"bin" successful'

clean_release:
	@rm -rf release/

release: clean_release
	@mkdir -p release/tmpl/src release/static
	@go build -ldflags="-s -w" -o release/wejay
	@$(INSTALL) static/* release/static/
	@$(INSTALL) tmpl/src/* release/tmpl/src/
	@$(INSTALL) .env-template release/.env-template
	@echo '"release" successful'
