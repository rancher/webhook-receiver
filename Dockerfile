FROM golang:1.12-alpine as builder
WORKDIR /go/src
ENV GO111MODULE on
COPY . .
RUN go build  -o /receiver -mod vendor cmd/main.go

FROM alpine as release
COPY --from=builder /receiver /receiver
ENTRYPOINT ["/receiver"]