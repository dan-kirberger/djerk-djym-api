set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -v -a -o release/linux/amd64/djerk-djym-api