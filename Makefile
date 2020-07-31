
define GetPluginVersion
$(shell node -p "'v' + require('./plugin.json').version")
endef

PLUGINVERSION=$(call GetPluginVersion)

dist: install-dependencies insertReleaseNotes quickdist removeReleaseNotes install-dependencies install-dependencies

define GetFromManifest
$(shell node -p "require('./plugin.json').$(1)")
endef

define InsertReleaseNotes
$(shell node -e
	"
	let fs = require('fs');
	try {
		let manifest = fs.readFileSync('plugin.json', 'utf8');
		manifest = JSON.parse(manifest);
		manifest.release_notes_url += manifest.version;
		let json = JSON.stringify(manifest, null, 2);
		fs.writeFileSync('plugin.json', json, 'utf8');
	} catch (err) {
		console.log(err);
	};"
)
endef

define RemoveReleaseNotes
$(shell node -e
	"
	let fs = require('fs');
	try {
		let manifest = fs.readFileSync('plugin.json', 'utf8');
		manifest = JSON.parse(manifest);
		if (manifest.release_notes_url.indexOf(manifest.version) >= 0) {
			manifest.release_notes_url = manifest.release_notes_url.substring(0, manifest.release_notes_url.indexOf(manifest.version));
		}
		let json = JSON.stringify(manifest, null, 2);
		fs.writeFileSync('plugin.json', json, 'utf8');
	} catch (err) {
		console.log(err);
	};"
)
endef

.PHONY: insertReleaseNotes removeReleaseNotes

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

clean:
	@echo Cleaning plugin

	rm -rf dist
	cd webapp && rm -rf node_modules
	cd webapp && rm -f .npminstall

install-dependencies:
	cd server && go mod download
	cd webapp && npm install

check-style: install-dependencies
	@echo Checking for style guide compliance

	@# TODO: configure lint for webapp
	@# cd webapp && npm run lint
	@# cd webapp && npm run check-types
	golangci-lint run ./...

release: dist
	@echo "Installing ghr"
	@go get -u github.com/tcnksm/ghr
	@echo "Create new tag"
	$(shell git tag $(PLUGINVERSION))
	@echo "Uploading artifacts"
	@ghr -t $(GITHUB_TOKEN) -u $(ORG_NAME) -r $(REPO_NAME) $(PLUGINVERSION) dist/
