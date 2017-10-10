.PHONY: assets cli

assets:
	(cd ui && yarn build)
	go-bindata -o ui/assets.go -pkg ui -prefix ui/build ui/build

cli: assets
	go build github.com/hayeah/qtum-portal/cli/qtumportal