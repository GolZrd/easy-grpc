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

test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/GolZrd/easy-grpc/internal/service/...,github.com/GolZrd/easy-grpc/internal/api/... -count 5

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/GolZrd/easy-grpc/internal/service/...,github.com/GolZrd/easy-grpc/internal/api/...
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore
