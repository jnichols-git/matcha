FORCE:

test: FORCE
	go test -coverprofile cicd/cover.out ./...

cover: test
	go tool cover -html cicd/cover.out

bench: FORCE
	go test -benchmem -bench=. ./...
