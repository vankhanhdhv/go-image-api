
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN go mod tidy
RUN go build -o main ./app.go
RUN ls -al

FROM scratch
WORKDIR /app
ARG DOCKER_METADATA_OUTPUT_VERSION=latest
ENV VERSION=$DOCKER_METADATA_OUTPUT_VERSION
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /app/go-image-api
COPY --from=builder /app/fonts /app/fonts
COPY --from=builder /app/index.html /app/index.html
EXPOSE 8080
ENTRYPOINT [ "/app/go-image-api" ]
