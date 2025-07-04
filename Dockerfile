FROM docker.io/library/golang:1.23-alpine AS builder
ARG GO_PRIVATE
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o core .

FROM docker.io/library/alpine:3.20
LABEL description="https://oxmix.net"
ARG VERSION
ENV VERSION=$VERSION
ARG VERSION_HASH
ENV VERSION_HASH=$VERSION_HASH
RUN apk --no-cache add sysbench fio speedtest-cli
COPY --from=builder /app/core .
ENTRYPOINT ["./core"]
