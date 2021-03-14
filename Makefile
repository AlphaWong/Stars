test:
		go test -count 1 -cover ./...
cover:
		go test -count 1 ./... -covermode=count -coverprofile=count.out && go tool cover -func=count.out && rm ./count.out
codecov:
		go test -count 1 ./... -covermode=count -coverprofile=count.out && go tool cover -func=count.out