.Phony: test

DATABASE_URL ?= postgresql://postgres:@localhost:5432/postgres?sslmode=disable
test:
	@DATABASE_URL=$(DATABASE_URL) go test -v
