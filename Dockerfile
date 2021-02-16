FROM scratch
ADD ./config_sample.json /config.json
COPY ./bin/api /
CMD ["/api"]
