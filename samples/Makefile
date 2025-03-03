GOPATH := $(shell go env GOPATH)

# USE K2 directly in real world
plan:	
	go run ../main.go plan  --inventory ./k2.inventory.yaml

apply:
	go run ../main.go apply --inventory ./k2.inventory.yaml

destroy:
	go run ../main.go destroy --inventory ./k2.inventory.yaml