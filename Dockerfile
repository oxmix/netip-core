FROM docker.io/library/golang:1.26-alpine AS builder
ARG GO_PRIVATE
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
# -ldflags="-s" removing symbol table (func name)
# -ldflags="-w" removing DWARF debug info (bad for: gdb, gcore, dlv)
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o core .
# for deep debug
#RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o core .

FROM docker.io/library/alpine:3.23
LABEL description="https://cloudnetip.com/wiki"
ARG VERSION
ENV VERSION=$VERSION
ARG VERSION_HASH
ENV VERSION_HASH=$VERSION_HASH
# who logged
RUN apk --no-cache add dbus elogind
# benchmarks
RUN apk --no-cache add sysbench fio speedtest-cli
# device metrics
RUN apk --no-cache add gcompat mdadm smartmontools zfs
COPY --from=builder /app/core .
ENTRYPOINT ["./core"]
