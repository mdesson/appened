
FROM --platform=linux/arm/v7 golang:latest AS builder
WORKDIR /go/src/github.com/mdesson/appended/
COPY ./clients ./clients
COPY ./HTTPLogger ./HTTPLogger
COPY ./go.mod .
COPY ./go.sum .
RUN go mod tidy
RUN cd clients/twilio && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM --platform=linux/arm/v7 alpine:latest  
WORKDIR /root/
RUN mkdir ../data/
COPY --from=builder /go/src/github.com/mdesson/appended/clients/twilio/app .
COPY --from=builder /go/src/github.com/mdesson/appended/clients/twilio/config.json .
CMD ["./app"]  
