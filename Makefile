

test:
	go test -coverprofile cicd/cover.out ./...

cover: test
	go tool cover -html cicd/cover.out
