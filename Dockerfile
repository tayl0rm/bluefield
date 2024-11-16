FROM golang:1.23

WORKDIR /src

COPY go/ /src

RUN CGO_ENABLED=0 GOOS=linux go build -o /bluefield cmd/main.go

CMD ["/bluefield"]