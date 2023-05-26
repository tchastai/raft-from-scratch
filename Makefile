build: 
	go build -o ./bin/raft-from-scratch

run: build
	./bin/raft-from-scratch

set_request:
	curl -X POST localhost:45000/set/toto-test

get_request: 
	curl -X GET localhost:45000/get/toto