test:
	go test ./pkg/...

test_cover:
	mkdir -p _testcover
	go test -v -coverprofile ./_testcover/cover.out ./pkg/...
	go tool cover -html=./_testcover/cover.out -o ./_testcover/cover.html
