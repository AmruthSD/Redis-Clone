OBJFILE := $(TEST:.go=)

build_main:
	go build -o ./bin/redis ./cmd/Redis

build_test:
	go build -o ./bin/$(OBJFILE) ./tests/$(TEST)

run_main:
	./bin/redis

run_test:
	./bin/$(OBJFILE)

clean:
	rm -f redis
