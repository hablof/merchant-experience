### генерация моков
mock-service:
	minimock -i ./internal/service.* -o ./internal/service
mock-router:
	minimock -i ./internal/router.* -o ./internal/router


### юнит-тесты
unit-test:
	go test -test.short ./...


### подъём базы в докере для тестов
up-test-db:
	docker run -p 5432:5432 -e POSTGRES_PASSWORD=1234 -e POSTGRES_DB=integration_testing -d --name testing-postgres --rm postgres

down-test-db:
	docker stop testing-postgres


### тесты с базой
database-test:
	go test internal/repository/repository_integration_test.go internal/repository/repository.go

integration-test:
	go test internal/*_test.go 


### запуски
run: 
	go run cmd/app/main.go

run-testserver:
	go run cmd/test-server/main.go

run-docker:


### подъём базы в докере для локальной работы
up-local-db:
	docker run -p 5432:5432 -e POSTGRES_PASSWORD=1234 -e POSTGRES_DB=merchant_experience -d --name local-postgres --rm postgres -v 

down-local-db:
	docker stop local-postgres
migrate-up:
	goose [OPTIONS] postgres "user=postgres password=1234 dbname=postgres sslmode=disable" COMMAND