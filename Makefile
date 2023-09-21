FORCE:

test:
	go test -coverprofile cicd/cover.out ./pkg/...

cover: test
	go tool cover -html cicd/cover.out

bench: FORCE
	go test -benchmem -bench=. ./pkg/...

license:
	addlicense -c "Decent Platforms" -l apache -y 2023 .
