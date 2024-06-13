obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver data_receiver/main.go
	@./bin/receiver

calculator:
	@go build -o bin/distance_calc distance_calc/main.go
	@./bin/distance_calc

invoice:
	@go build -o bin/invoicer invoicer/main.go
	@./bin/invoicer

.PHONY: obu