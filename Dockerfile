FROM golang:1.14 as base
COPY main.go Makefile go.mod go.sum /code/
COPY cmd /code/cmd/
COPY pkg /code/pkg/
WORKDIR /code
RUN make build

FROM scratch
COPY --from=base /code/bin/vvalidator-linux /vvalidator-linux
ENTRYPOINT ["/vvalidator-linux"]