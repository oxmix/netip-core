FROM docker.io/library/golang:1.22-alpine AS builder
ARG GO_PRIVATE
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o core .

FROM docker.io/library/alpine:3.18
LABEL description="https://oxmix.net"
ARG VERSION
ENV VERSION=$VERSION
ARG VERSION_HASH
ENV VERSION_HASH=$VERSION_HASH
RUN apk --no-cache add sysbench fio speedtest-cli
COPY --from=builder /app/core .
ENTRYPOINT ["./core"]
