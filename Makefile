export VERSION

local: test config.json
	@if [ ! -z $$(pgrep keystore) ]; then \
	  echo "killing running process" ;\
	  kill $$(pgrep keystore) ;\
	fi ;\
	go-bindata data/
	go build -o bin/keystore-local
	bin/keystore-local --config config.json &
	@echo "keystore is viewable at http://localhost:<port>"

config.json:
	@echo "no config.json found, running generateConfig.py"
	./generateConfig.py config.json

test:
	go test -v keystore/...

deploy: test build config.production.json infrastructure
	test -n "$(VERSION)" # $$VERSION
	$(MAKE) -C infrastructure $@

build:
	go-bindata data/
	go build -o bin/keystore -ldflags="-s -w"

config.production.json:
	@echo "no config.production.json found, running generateConfig.py"
	./generateConfig.py config.production.json

clean:
	rm -rf bin/*

.PHONY: infrastructure clean
