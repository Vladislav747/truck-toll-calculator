obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver data_receiver/main.go
	@./bin/receiver

calculator:
	@go build -o bin/distance_calc distance_calc/main.go
	@./bin/distance_calc

invoicer:
	@go build -o bin/invoicer invoicer/main.go
	@./bin/invoicer

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto

proto_grpc:
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto
.PHONY: obu invoicer