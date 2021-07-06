run:
	go run ./cmd/api

db/psql:
	psql ${GREENLIGHT_DB_DSN}

db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up

db/migrations/new:
	@echo "Creating migrations for ${name}"
	migrate create -seq -ext=.sql -dir-./migrations ${name}

