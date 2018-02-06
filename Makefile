export VERSION

local: test config.json front-end
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

#TODO this npm/node stuff should probably be done better
#only build stale files etc
front-end:
	$(MAKE) -C web $@

deploy: test build config.production.json infrastructure
	test -n "$(VERSION)" # $$VERSION
	$(MAKE) -C infrastructure $@

js:
	$(MAKE) -C web $@

css:
	$(MAKE) -C web $@

build: front-end
	go-bindata data/
	go build -o bin/keystore -ldflags="-s -w"

config.production.json:
	@echo "no config.production.json found, running generateConfig.py"
	./generateConfig.py config.production.json

clean:
	rm -rf bin/*

.PHONY: infrastructure clean
