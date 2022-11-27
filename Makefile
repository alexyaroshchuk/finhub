run: up migrateup

build: up_with_build migrateup

up:
	docker-compose -f docker-compose.yml up -d

up_with_build:
	docker-compose -f docker-compose.yml up --build -d

migrateup:
	migrate -path db/migration -database "postgresql://test:password@localhost:5432/finhub?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://<user>@<pwd>:5432/finhub?sslmode=disable" -verbose down

test:
	./test.sh