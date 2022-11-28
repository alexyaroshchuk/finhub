run:
	docker compose up -d

build:
	docker compose up --build -d

migrateup:
	migrate -path db/migration -database "postgresql://test:password@localhost:5432/finhub?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://test:password@localhost:5432/finhub?sslmode=disable" -verbose down

test:
	./test.sh