gen:
	protoc --go_out=pb \
    --proto_path=proto proto/*.proto

clean:
	rm pb/*.go

run:
	go run main.go