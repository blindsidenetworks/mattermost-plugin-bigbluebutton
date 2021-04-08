define GetPluginId
$(shell node -p "require('./plugin.json').id")
endef

PLUGINNAME=$(call GetPluginId)

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

.PHONY: deploy
.SILENT:
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
