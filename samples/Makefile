GOPATH := $(shell go env GOPATH)

# USE K2 directly in real world
render-plan:	
	go run ../main.go render-plan  --inventory ./k2.inventory.yaml

render:
	go run ../main.go render --inventory ./k2.inventory.yaml

unrender:
	go run ../main.go unrender --inventory ./k2.inventory.yaml


stack-up: 
	go run ../main.go stack up stack-one 


stack-down:
	go run ../main.go stack down stack-one 