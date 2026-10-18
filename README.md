# Сервис «Рестораны» (вариант 5), Go kit: HTTP + gRPC

## Запуск
    go mod tidy        # скачает зависимости, создаст go.sum
    go run .

HTTP: :8080, gRPC: :9090 (меняется флагами -http.addr и -grpc.addr).

## Проверка HTTP (curl)
    curl -X POST http://localhost:8080/restaurants -H "Content-Type: application/json" \
      -d '{"name":"Пицца-Хаус","address":"ул. Ленина, 1","rating":4.2,"open":true}'
    curl http://localhost:8080/restaurants/<id>
    curl -X PATCH http://localhost:8080/restaurants/<id>/rating -H "Content-Type: application/json" \
      -d '{"rating":4.9}'

В Windows PowerShell лучше использовать curl.exe и экранировать кавычки, либо Postman.

## Проверка gRPC (grpcurl, reflection включён)
    grpcurl -plaintext localhost:9090 list
    grpcurl -plaintext -d '{"restaurant":{"name":"Суши","address":"пр. Мира, 5","rating":4,"open":true}}' localhost:9090 restaurant.RestaurantService/Create
    grpcurl -plaintext -d '{"id":"<id>"}' localhost:9090 restaurant.RestaurantService/GetByID
    grpcurl -plaintext -d '{"id":"<id>","rating":4.5}' localhost:9090 restaurant.RestaurantService/ChangeRating

## Перегенерация protobuf (не обязательно, pb/ уже сгенерирован)
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    protoc -I proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/restaurant.proto
