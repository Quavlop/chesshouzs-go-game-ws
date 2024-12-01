FROM golang:1.23.1

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go mod verify

COPY . .

RUN ln -sf /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && echo "Asia/Jakarta" > /etc/timezone

RUN CGO_ENABLED=1 GOOS=linux go build -o /docker-gs-ping

CMD ["/docker-gs-ping"]