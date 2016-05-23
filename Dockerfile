FROM golang:latest

RUN go get github.com/constabulary/gb/...

RUN mkdir -p /opt/oden
RUN curl -L https://github.com/oden-lang/oden/releases/download/0.3.1/oden-0.3.1-linux.tar.gz | tar -xvz -C /opt

RUN mkdir -p /app
WORKDIR /app

COPY . /app
RUN gb build

EXPOSE 8080
ENV GOPATH=/go
ENV ODEN_CLI=/opt/oden/bin/oden

ENTRYPOINT ["bin/playground"]
