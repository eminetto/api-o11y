.PHONY: all
all: build
FORCE: ;

.PHONY: build

build: build-auth build-feedback build-vote

build-auth:
	cd auth; go build -o bin/auth main.go

build-feedback:
	cd feedbacks; go build -o bin/feedbacks main.go

build-vote:
	cd votes; go build -o bin/votes main.go

build-auth-linux:
	cd auth; CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags "netgo" -installsuffix netgo -o bin/auth main.go

build-feedbacks-linux:
	cd feedbacks; CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags "netgo" -installsuffix netgo -o bin/feedbacks main.go

build-votes-linux:
	cd votes; CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags "netgo" -installsuffix netgo -o bin/votes main.go

build-linux: build-auth-linux build-feedbacks-linux build-votes-linux

build-auth-docker: build-auth-linux
	cd auth; docker build -t auth -f Dockerfile .

build-feedbacks-docker: build-feedbacks-linux
	cd feedbacks; docker build -t feedbacks -f Dockerfile .

build-votes-docker: build-votes-linux
	cd votes; docker build -t votes -f Dockerfile .

run-docker: build-auth-docker build-feedbacks-docker build-votes-docker
    docker run -d -p 8081:8081 auth
    docker run -d -p 8082:8082 feedbacks
    docker run -d -p 8083:8083 votes

