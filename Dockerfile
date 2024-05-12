FROM golang:1.11.1 as builder
WORKDIR /build
ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make static-build

FROM scratch
WORKDIR /root/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /build/bot .
COPY --from=builder /build/pokemon.csv .

EXPOSE 8080

CMD ["./bot"]  