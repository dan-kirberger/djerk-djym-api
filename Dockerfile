FROM plugins/base:multiarch

ADD release/linux/amd64/djerk-djym-api /bin/

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "/bin/djerk-djym-api", "-ping" ]

ENTRYPOINT ["/bin/helloworld"]