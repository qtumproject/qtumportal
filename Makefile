.PHONY: assets

assets:
	(cd ui && yarn build)
	go-bindata -o ui/assets.go -pkg ui -prefix ui/build  ui/build