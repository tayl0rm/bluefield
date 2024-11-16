FROM golang:1.23

WORKDIR /bluefield

COPY go/* ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build cmd/main.go -o .

CMD ["/bluefield"]