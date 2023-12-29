FROM golang:1.21 AS builder

LABEL maintainer = "Giwa Oluwatobi, giwaoluwatobi@gmail.com"

WORKDIR /app

COPY . /app

RUN go mod tidy
RUN go build -o bin/BlockRide cmd/blockride/main.go

FROM alpine 

WORKDIR /root/

COPY --from=builder /app/bin/BlockRide /root/

CMD ./BlockRide