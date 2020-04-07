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
	@rm -rf release/wejay

release: clean_release
	@mkdir -p release/wejay/tmpl/src release/wejay/static
	@go build -ldflags="-s -w" -o release/wejay/bin
	@$(INSTALL) static/* release/wejay/static/
	@$(INSTALL) tmpl/src/* release/wejay/tmpl/src/
	@$(INSTALL) .env-template release/wejay/.env-template
	@echo '"release" successful'
