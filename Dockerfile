FROM plugins/base:multiarch

ADD release/linux/amd64/djerk-djym-api /bin/

EXPOSE 8080

ENTRYPOINT ["/bin/djerk-djym-api"]