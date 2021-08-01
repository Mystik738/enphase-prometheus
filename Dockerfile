FROM golang AS build

WORKDIR /go/src/app
COPY . .

RUN go env -w GO111MODULE=auto
RUN go env -w CGO_ENABLED=0
RUN go get -d -v ./...
RUN go build -v enphase-prometheus.go

FROM scratch
COPY --from=build /go/src/app/enphase-prometheus app

CMD ["./app"]