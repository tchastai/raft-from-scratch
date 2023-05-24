build: 
	go build -o ./bin/raft-from-scratch

run: build
	./bin/raft-from-scratch

test:
	go test ./...