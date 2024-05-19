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

- Запускаем consumer
```
make calculator
```