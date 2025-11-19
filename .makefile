.PHONY: build run test migrate clean

build:
	docker-compose build

run:
	docker-compose up

test:
	go test ./...

migrate:
	docker-compose run migration

clean:
	docker-compose down -v
	rm -f main

db-shell:
	docker-compose exec db psql -U postgres -d qna_db