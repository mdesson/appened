
FROM golang:latest AS builder
WORKDIR /go/src/github.com/mdesson/appended/clients/twilio
COPY ./cmd/ ./cmd/
COPY ./note/ ./note/
COPY ./HTTPLogger/ ./HTTPLogger/
COPY ./HTTPLogger/ ./HTTPLogger/
COPY ./go.mod .
COPY ./go.sum .
RUN cd cmd && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
RUN ls cmd

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
RUN mkdir ../data/
COPY --from=builder /go/src/github.com/mdesson/appended/cmd/clients/twilio/app .
COPY auth.sh .
RUN ./auth.sh
CMD ["./app"]  
