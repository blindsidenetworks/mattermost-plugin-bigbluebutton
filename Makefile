.PHONY: dist build quickdist quickbuild install-dependencies clean check-style check-style-server golangci-lint test insertReleaseNotes removeReleaseNotes deploy release

define GetPluginVersion
$(shell node -p "'v' + require('./plugin.json').version")
endef

PLUGINVERSION=$(call GetPluginVersion)

define GetPluginId
$(shell node -p "require('./plugin.json').id")
endef

PLUGINNAME=$(call GetPluginId)

dist: install-dependencies
	$(MAKE) insertReleaseNotes
	$(MAKE) quickdist; result=$$?; $(MAKE) removeReleaseNotes; exit $$result

build: install-dependencies
	$(MAKE) quickbuild

export GOFLAGS=-trimpath

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

# Build directly with the Go toolchain; GOPATH-mode go get and gox are obsolete.
quickdist: quickbuild

quickbuild:
	mkdir -p dist
	cd webapp && npm run build
	@set -e; for platform in darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64; do \
		os=$${platform%/*}; arch=$${platform#*/}; \
		stage=$$(mktemp -d); \
		mkdir -p "$$stage/bigbluebutton/server" "$$stage/bigbluebutton/webapp"; \
		cp plugin.json "$$stage/bigbluebutton/"; \
		cp -R assets "$$stage/bigbluebutton/"; \
		cp -R webapp/dist/. "$$stage/bigbluebutton/webapp/"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -ldflags="-X main.PluginVersion=$(call GetFromManifest,version)" -o "$$stage/bigbluebutton/server/plugin.exe" ./server; \
		tar -czf "dist/bigbluebutton_$${os}_$${arch}.tar.gz" -C "$$stage" bigbluebutton; \
		rm -rf "$$stage"; \
	done

install-dependencies:
	go mod download
	cd webapp && npm ci

clean:
	@echo Cleaning generated plugin bundles
	rm -rf dist webapp/dist

check-style: check-style-server
	cd webapp && npm run check

golangci-lint:
	golangci-lint run ./server/...

check-style-server: golangci-lint
	go vet ./server/...
	@test -z "$$(gofmt -l server)" || (gofmt -l server; exit 1)

test:
	go test -race ./server/...
	cd webapp && npm test

release: dist
	@echo "Installing ghr"
	@go install github.com/tcnksm/ghr@v0.18.6
	@echo "Create new tag"
	$(shell git tag $(PLUGINVERSION))
	@echo "Uploading artifacts"
	@ghr -t $(GITHUB_TOKEN) -u $(ORG_NAME) -r $(REPO_NAME) $(PLUGINVERSION) dist/

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

