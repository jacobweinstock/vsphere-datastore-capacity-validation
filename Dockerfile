FROM golang:1.14 as base
COPY . /code
WORKDIR /code
RUN make build

FROM scratch
COPY --from=base /code/bin/vvalidator-linux /vvalidator-linux
ENTRYPOINT ["/vvalidator-linux"]