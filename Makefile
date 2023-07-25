mock-service:
	minimock -i ./internal/service.* -o ./internal/service
mock-router:
	minimock -i ./internal/router.* -o ./internal/router


unit-test:
	go test -test.short ./...


# тесты с реальной базой
up-test-db:
	docker run -p 5432:5432 -e POSTGRES_PASSWORD=1234 -e POSTGRES_DB=integration_testing -d --name testing-postgres --rm postgres

db-test:
	go test internal/repository/repository_integration_test.go internal/repository/repository.go

down-test-db:
	docker stop testing-postgres