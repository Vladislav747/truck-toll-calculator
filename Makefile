obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver data_receiver/main.go
	@./bin/receiver

calculator:
	@go build -o bin/distance_calc distance_calc/main.go
	@./bin/distance_calc

aggregator:
	@go build -o bin/aggregator aggregator/main.go
	@./bin/aggregator

temp_temp:
	@go build -o bin/temp_code temp/main.go
	@./bin/temp_code

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto

proto_grpc:
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto
.PHONY: obu invoicer