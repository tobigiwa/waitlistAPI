FROM golang:1.21 AS builder

LABEL maintainer = "Giwa Oluwatobi, giwaoluwatobi@gmail.com"

WORKDIR /app

COPY go.mod go.sum ./

# vendoring instead
# RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o /app/blockride cmd/blockride/main.go

FROM alpine 

WORKDIR /

COPY --from=builder /app/blockride /blockride

RUN chmod +x blockride

ENTRYPOINT [ "/blockride" ]