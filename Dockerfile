FROM as4go

MAINTAINER Joern Weissenborn <joern.weissenborn@gmail.com>

RUN go get github.com/joernweissenborn/propertygraph2go
RUN go get code.google.com/p/go.net/websocket

EXPOSE 5555 5557


COPY . /var/local/gopath/src/aursirrt/
WORKDIR /var/local/gopath/src/aursirrt/
RUN go build src/main.go


ENTRYPOINT ["sh","dockerinit.sh"]

