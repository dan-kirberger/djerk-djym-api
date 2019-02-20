FROM golang:1.11.5

ADD release/linux/amd64/djerk-djym-api /bin/

EXPOSE 8080

RUN chmod +x /bin/djerk-djym-api

ENTRYPOINT ["/bin/djerk-djym-api"]