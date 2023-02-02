ENV_NAME := $(shell git rev-parse --abbrev-ref HEAD)
nuke:
	@npx lerna run nuke
	@npx lerna clean -y
	@npm run nuke
	@npm install
	@npx lerna bootstrap

install:
	@npm install
	@npx lerna bootstrap
	@npx @tuxounet-k2/cli apply

clean:
	@npx @tuxounet-k2/cli clean

build: clean
	@npx lerna exec npm run build
	@npx lerna bootstrap

test:
	@npx lerna run test

deploy: install build test
	@K2_ENV=${ENV_NAME} npx lerna run deploy

undeploy: install build
	@K2_ENV=${ENV_NAME} npx lerna run synth
	@K2_ENV=${ENV_NAME} npx lerna run destroy


apply:
	@npx @tuxounet-k2/cli apply
	@npx lerna bootstrap


list:
	@npx @tuxounet-k2/cli list