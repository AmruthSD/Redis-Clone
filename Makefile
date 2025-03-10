build:
	go build -o redis ./cmd/Redis

run:
	./redis

clean:
	rm -f redis
