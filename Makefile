build: install-dependencies quickbuild

define GetFromManifest
$(shell node -p "require('./plugin.json').$(1)")
endef

quickbuild:
	@echo Building plugin

	rm -rf dist
	cd server && go get github.com/mitchellh/gox
	$(shell go env GOPATH)/bin/gox -ldflags="-X main.PluginVersion=$(call GetFromManifest,version)" -osarch='darwin/amd64 linux/amd64 windows/amd64' -gcflags='all=-N -l' -output 'dist/intermediate/plugin_{{.OS}}_{{.Arch}}' ./server

	mkdir -p dist/bigbluebutton/server

	# Clean old dist
	rm -rf webapp/dist
	cd webapp && npm run build

	# Copy files from webapp
	mkdir -p dist/bigbluebutton/webapp
	cp webapp/dist/* dist/bigbluebutton/webapp/

	# Copy plugin files
	cp plugin.json dist/bigbluebutton/

	# Package darwin pakckage
	mv dist/intermediate/plugin_darwin_amd64 dist/bigbluebutton/server/plugin.exe
	cd dist && tar -zcvf bigbluebutton_darwin_amd64.tar.gz bigbluebutton/*

	# Package linux package
	mv dist/intermediate/plugin_linux_amd64 dist/bigbluebutton/server/plugin.exe
	cd dist && tar -zcvf bigbluebutton_linux_amd64.tar.gz bigbluebutton/*

	# Package windows package
	mv dist/intermediate/plugin_windows_amd64.exe dist/bigbluebutton/server/plugin.exe
	cd dist && tar -zcvf bigbluebutton_windows_amd64.tar.gz bigbluebutton/*

	# Clean up temp files
	rm -rf dist/bigbluebutton
	rm -rf dist/intermediate

	@echo Plugin built at: dist/bigbluebutton.tar.gz

install-dependencies:
	cd server && go get github.com/Masterminds/glide
	cd server && $(shell go env GOPATH)/bin/glide install

	#installs node modules
	cd webapp && npm install

clean:
	@echo Cleaning plugin

	rm -rf dist
	cd webapp && rm -rf node_modules
	cd webapp && rm -f .npminstall

check-style: check-style-server

check-style-server:
	@echo Running GOFMT

	@for package in $$(go list ./server/...); do \
		echo "Checking "$$package; \
		files=$$(go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' $$package); \
		if [ "$$files" ]; then \
			gofmt_output=$$(gofmt -w -s $$files 2>&1); \
			if [ "$$gofmt_output" ]; then \
				echo "$$gofmt_output"; \
				echo "gofmt failure"; \
				exit 1; \
			fi; \
		fi; \
	done
	@echo "gofmt success"; \
