FROM golang:1.21 AS builder

LABEL maintainer = "Giwa Oluwatobi, giwaoluwatobi@gmail.com"

WORKDIR /app

COPY go.mod go.sum ./

# vendoring instead
# RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o /app/companyXYZ cmd/companyXYZ/main.go

FROM alpine 

WORKDIR /

COPY --from=builder /app/companyXYZ /companyXYZ

RUN chmod +x companyXYZ

ENTRYPOINT [ "/companyXYZ" ]