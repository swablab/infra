FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:edge
RUN apk add --no-cache font-ubuntu font-noto typst git
COPY --from=builder /go/src/app/typst-api /typst-api
EXPOSE 8080
ENTRYPOINT ["/typst-api"]
