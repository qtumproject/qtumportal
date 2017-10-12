.PHONY: assets cli

assets:
	(cd ui && yarn build)
	go-bindata -o ui/assets.go -pkg ui -prefix ui/build ui/build

cli:
	GOOS=darwin go build -o qtumportal-mac github.com/hayeah/qtum-portal/cli/qtumportal
	GOOS=linux go build -o qtumportal-linux github.com/hayeah/qtum-portal/cli/qtumportal