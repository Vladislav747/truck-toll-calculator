# toll-calculator
```
docker run --name kafka -p 9092:9092 -e ALLOW_PLAINTEXT_LISTENER=yes -e KAFKA_CFG_AUTO
```



- Сначала запускаем брокера через
```
docker-compose up -d
```

- Запускаем receiver(producer) который получает данные из ws и начинает их produce в kafka
```
make receiver
```

- Запускаем obu - генерирующий данный с websocket в receiver
```
make obu
```

- Запускаем invoice - куда мы обрабатываем данные через запрос - типо микросервиса

```
make invoicer
```

- Запускаем consumer
```
make calculator
```

## Installing GRPC and ProtoBuffer plugins for Golang
1. Protobuffers
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28   
```

2. GRPC

```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

3. NOTE that you need to set the /go/bin directory in your path
Just like this or whatever your go directly lives.

```
PATH="$PATH:${HOME}/go/bin"
export PATH=$GOPATH/bin:$PATH

```


4 Install the package dependencies
```
go get google.golang.org/protobuf
```

4.2 grpc package
```
go get google.golang.org/grpc/
```


5 Запустить proto
```
make proto
```