FROM plugins/base:multiarch

ADD release/linux/amd64/djerk-djym-api /bin/

CMD [ "/bin/djerk-djym-api", "-ping" ]

EXPOSE 8080

ENTRYPOINT ["/bin/djerk-djym-api"]