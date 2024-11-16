FROM golang:1.23

WORKDIR /app

COPY go/*.go ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build cmd/main.go -o /bluefield

CMD ["/bluefield"]