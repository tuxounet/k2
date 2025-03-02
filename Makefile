run: 
	go run ./main.go

build:
	go build -o ./.out/k2 ./main.go

plan:
	go run ./main.go plan  --inventory ./samples/k2.inventory.yaml

apply:
	go run ./main.go apply --inventory ./samples/k2.inventory.yaml

destroy:
	go run ./main.go destroy --inventory ./samples/k2.inventory.yaml