FROM golang:1.21 AS builder

LABEL maintainer = "Giwa Oluwatobi, giwaoluwatobi@gmail.com"

WORKDIR /app

COPY . /app

RUN make setup

FROM alpine 

WORKDIR /root/

COPY --from=builder /app/bin/ /root/

CMD /bin/BlockRide