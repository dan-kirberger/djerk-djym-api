build_linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -o release/linux/amd64/djerk-djym-api
docker:
	docker build -t dan-kirberger/djerk-djym-api .

test:
	go test -v .