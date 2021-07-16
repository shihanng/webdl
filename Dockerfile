# syntax=docker/dockerfile:1
FROM golang:1.16 as build
WORKDIR /go/src/github.com/shihanng/webdl/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o webdl .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY entrypoint.sh entrypoint.sh
COPY --from=build /go/src/github.com/shihanng/webdl/webdl .
ENTRYPOINT ["/app/entrypoint.sh"]
