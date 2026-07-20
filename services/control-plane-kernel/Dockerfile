FROM golang:1.26-alpine AS build
WORKDIR /src
COPY services/control-plane-kernel/go.mod ./
COPY services/control-plane-kernel/ ./
RUN /usr/local/go/bin/go test ./... && /usr/local/go/bin/go build -o /out/control-plane-kernel ./cmd/kernel

FROM alpine:3.22
RUN adduser -D -H fusion
USER fusion
COPY --from=build /out/control-plane-kernel /usr/local/bin/control-plane-kernel
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/control-plane-kernel"]
