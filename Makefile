build:
	@echo Building plugin

	rm -rf dist

	cd server && go get github.com/Masterminds/glide
	cd server && $(shell go env GOPATH)/bin/glide install

	cd server && go build -o plugin.exe plugin.go responsehandlers.go helpers.go config.go

	mkdir -p dist/bigbluebutton/server
	cp server/plugin.exe dist/bigbluebutton/server

	# Clean old dist
	rm -rf webapp/dist
	#installs node modules
	cd webapp && npm install
	cd webapp && npm run build


	# Copy files from webapp
	mkdir -p dist/bigbluebutton/webapp
	cp webapp/dist/* dist/bigbluebutton/webapp/

	# Copy plugin files
	cp plugin.yaml dist/bigbluebutton/

	# Compress
	cd dist && tar -zcvf bigbluebutton.tar.gz bigbluebutton/*

	# Clean up temp files
	rm -rf dist/bigbluebutton

	@echo Plugin built at: dist/bigbluebutton.tar.gz

clean:
	@echo Cleaning plugin

	rm -rf dist
	cd webapp && rm -rf node_modules
	cd webapp && rm -f .npminstall
