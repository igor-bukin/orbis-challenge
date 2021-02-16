FROM scratch
ADD ./config_sample.json /config.json
ADD ./src/ca-certificates.crt /etc/ssl/certs/
COPY ./bin/api /
COPY ./_swaggerui /_swaggerui
CMD ["/api"]
