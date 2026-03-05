GOPATH := $(shell go env GOPATH)

# USE K2 directly in real world
render-plan:	
	go run ../main.go render-plan  --inventory ./k2.inventory.yaml

render:
	go run ../main.go render --inventory ./k2.inventory.yaml

unrender:
	go run ../main.go unrender --inventory ./k2.inventory.yaml