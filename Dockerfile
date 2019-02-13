FROM plugins/base:multiarch

ADD release/linux/amd64/helloworld /bin/

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "/bin/djerk-djym-api", "-ping" ]

ENTRYPOINT ["/bin/helloworld"]