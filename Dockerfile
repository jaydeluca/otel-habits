# Start from golang base image
FROM golang:1.21-alpine3.18 as builder

# Enable go modules
ENV GO111MODULE=on

# Install git. (alpine image does not have git in it, needed for gcc)
RUN apk update && apk add --no-cache git && apk add bash && apk add build-base

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY ./ .

# Build the application.
RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o ./main .


# Run executable
CMD ["./main"]