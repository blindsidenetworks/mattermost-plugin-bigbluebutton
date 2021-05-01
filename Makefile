.PHONY: insertReleaseNotes removeReleaseNotes deploy

define GetPluginVersion
$(shell node -p "'v' + require('./plugin.json').version")
endef

PLUGINVERSION=$(call GetPluginVersion)

define GetPluginId
$(shell node -p "require('./plugin.json').id")
endef

PLUGINNAME=$(call GetPluginId)

dist: install-dependencies insertReleaseNotes quickdist removeReleaseNotes install-dependencies install-dependencies

build: install-dependencies quickbuild

define GetFromManifest
$(shell node -p "require('./plugin.json').$(1)")
endef

define InsertReleaseNotes
$(shell node -e \
	"\
	let fs = require('fs');\
	try {\
		let manifest = fs.readFileSync('plugin.json', 'utf8');\
		manifest = JSON.parse(manifest);\
		manifest.release_notes_url += manifest.version;\
		let json = JSON.stringify(manifest, null, 2);\
		fs.writeFileSync('plugin.json', json, 'utf8');\
	} catch (err) {\
		console.log(err);\
	};"\
)
endef

define RemoveReleaseNotes
$(shell node -e\
	"\
	let fs = require('fs');\
	try {\
		let manifest = fs.readFileSync('plugin.json', 'utf8');\
		manifest = JSON.parse(manifest);\
		if (manifest.release_notes_url.indexOf(manifest.version) >= 0) {\
			manifest.release_notes_url = manifest.release_notes_url.substring(0, manifest.release_notes_url.indexOf(manifest.version));\
		}\
		let json = JSON.stringify(manifest, null, 2);\
		fs.writeFileSync('plugin.json', json, 'utf8');\
	} catch (err) {\
		console.log(err);\
	};"\
)
endef

insertReleaseNotes:
	$(call InsertReleaseNotes)

removeReleaseNotes:
	$(call RemoveReleaseNotes)

quickdist:
	@echo Building plugin

	rm -rf dist
	cd server && go get github.com/mitchellh/gox
	$(shell go env GOPATH)/bin/gox -ldflags="-X main.PluginVersion=$(call GetFromManifest,version)" -osarch='darwin/amd64 linux/amd64 windows/amd64' -gcflags='all=-N -l' -output 'dist/intermediate/plugin_{{.OS}}_{{.Arch}}' ./server

	mkdir -p dist/bigbluebutton/server

quickbuild:
	rm -rf dist/
	mkdir -p dist/bigbluebutton
	cp plugin.json dist/bigbluebutton
	cp -r assets dist/bigbluebutton
	cd webapp && npm run build
	mkdir -p dist/bigbluebutton/webapp
	cp -r webapp/dist/ dist/bigbluebutton/webapp
	cd server && GO111MODULE=off go get github.com/mitchellh/gox
	$(shell go env GOPATH)/bin/gox -ldflags="-X main.PluginVersion=$(call GetFromManifest,version)" -osarch='darwin/amd64 linux/amd64 windows/amd64' -gcflags='all=-N -l' -output 'dist/intermediate/plugin_{{.OS}}_{{.Arch}}' ./server
	mkdir -p dist/bigbluebutton/server

	cp dist/intermediate/plugin_darwin_amd64 dist/bigbluebutton/server/plugin.exe
	cd dist && tar -cvzf bigbluebutton_darwin_amd64.tar.gz bigbluebutton

	cp dist/intermediate/plugin_linux_amd64 dist/bigbluebutton/server/plugin.exe
	cd dist && tar -cvzf bigbluebutton_linux_amd64.tar.gz bigbluebutton

	cp dist/intermediate/plugin_windows_amd64.exe dist/bigbluebutton/server/plugin.exe
	cd dist && tar -cvzf bigbluebutton_windows_amd64.tar.gz bigbluebutton

	rm -rf dist/bigbluebutton
	rm -rf dist/intermediate

	@echo Plugin built at: dist/bigbluebutton.tar.gz

install-dependencies: clean
	go mod tidy
	go mod vendor

	#installs node modules
	cd webapp && npm install

clean:
	@echo Cleaning plugin

	rm -rf dist
	rm -rf vendor
	cd webapp && rm -rf node_modules
	cd webapp && rm -f .npminstall

check-style: install-dependencies check-style-server
	@echo Checking for style guide compliance

	@# TODO: configure lint for webapp
	@# cd webapp && npm run lint
	@# cd webapp && npm run check-types


release: dist
	@echo "Installing ghr"
	@go get -u github.com/tcnksm/ghr
	@echo "Create new tag"
	$(shell git tag $(PLUGINVERSION))
	@echo "Uploading artifacts"
	@ghr -t $(GITHUB_TOKEN) -u $(ORG_NAME) -r $(REPO_NAME) $(PLUGINVERSION) dist/

golangci-lint:
	golangci-lint run ./...

check-style-server: golangci-lint
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

deploy:
	echo "Installing plugin via API"

	echo "Authenticating admin user..." && \
	TOKEN=`http --print h POST $(MM_SERVICESETTINGS_SITEURL)/api/v4/users/login login_id=$(MM_ADMIN_USERNAME) password=$(MM_ADMIN_PASSWORD) X-Requested-With:"XMLHttpRequest" | grep Token | cut -f2 -d' '` && \
	http GET $(MM_SERVICESETTINGS_SITEURL)/api/v4/users/me Authorization:"Bearer $$TOKEN" > /dev/null && \
	echo "Deleting existing plugin..." && \
	http DELETE $(MM_SERVICESETTINGS_SITEURL)/api/v4/plugins/$(PLUGINNAME) Authorization:"Bearer $$TOKEN" > /dev/null && \
	echo "Uploading plugin..." && \
	http --check-status --form POST $(MM_SERVICESETTINGS_SITEURL)/api/v4/plugins plugin@dist/$(PLUGINNAME)_$(PLATFORM)_amd64.tar.gz Authorization:"Bearer $$TOKEN" > /dev/null && \
	echo "Enabling uploaded plugin..." && \
	http POST $(MM_SERVICESETTINGS_SITEURL)/api/v4/plugins/$(PLUGINNAME)/enable Authorization:"Bearer $$TOKEN" > /dev/null && \
	echo "Logging out admin user..." && \
	http POST $(MM_SERVICESETTINGS_SITEURL)/api/v4/users/logout Authorization:"Bearer $$TOKEN" > /dev/null && \
	echo "Plugin uploaded successfully"

