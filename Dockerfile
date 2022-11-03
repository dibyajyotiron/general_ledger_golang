FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR /
COPY . /
RUN go build .

EXPOSE 3000
ENTRYPOINT ["general_ledger_golang/cmd/http/main.go"]
