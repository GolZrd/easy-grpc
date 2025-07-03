include .env

local-migration-status:
	goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} status -v

local-migration-up:
	goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} up -v

local-migration-down:
	goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} down -v


generate:
	make generate-note-api

generate-note-api:
	mkdir pkg\note_v1
	protoc --proto_path api/note_v1 \
	--go_out=pkg/note_v1 --go_opt=paths=source_relative \
	--go-grpc_out=pkg/note_v1 --go-grpc_opt=paths=source_relative \
	api\note_v1\note.proto
