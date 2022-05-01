FROM golang:1.17 as builder

WORKDIR /app

COPY . .

RUN make go-install

FROM alpine

COPY --from=builder /go/bin/pinkPanther /bin/pinkPanther

EXPOSE 8080

CMD [ "/bin/pinkPanther" ]
