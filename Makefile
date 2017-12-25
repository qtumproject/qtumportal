.PHONY: assets cli

assets:
	(cd ui && yarn build)
	go-bindata -o ui/assets.go -pkg ui -prefix ui/build ui/build

abiplay-assets:
	# (cd ui && yarn build)
	mkdir -p assets/abiplay
	go-bindata -o assets/abiplay/assets.go -pkg abiplay -prefix abiplay/build abiplay/build
	# go-bindata -o assets/abiplay/assets.go -pkg abiplay abiplay/build

authuiui-assets:
	# (cd ui && yarn build)
	mkdir -p assets/authui
	go-bindata -o assets/authui/assets.go -pkg authui -prefix ui/build ui/build

cli:
	GOOS=darwin go build -o build/qtumportal-mac github.com/hayeah/qtum-portal/cli/qtumportal
	GOOS=linux go build -o build/qtumportal-linux github.com/hayeah/qtum-portal/cli/qtumportal