test:
		go test -count 1 -race -cover --covermode=atomic -coverpkg=./... ./...
test-rase:
		go test -count 1 -race -cover --covermode=atomic -coverpkg=./... ./...
cover:
		go test -count 1 -race -cover --covermode=atomic ./... -coverpkg=./... -coverprofile=count.out && go tool cover -func=count.out && rm ./count.out
codecov:
		go test -count 1 -race -cover --covermode=atomic ./... -coverpkg=./... -coverprofile=count.out && go tool cover -func=count.out