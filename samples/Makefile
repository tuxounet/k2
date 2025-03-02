GOPATH := $(shell go env GOPATH)
install: 
	go install github.com/tuxounet/k2@v0.9.1

plan:
	${GOPATH}/bin/k2 plan  --inventory ./k2.inventory.yaml

apply:
	${GOPATH}/bin/k2  apply --inventory ./k2.inventory.yaml

destroy:
	${GOPATH}/bin/k2  destroy --inventory ./k2.inventory.yaml