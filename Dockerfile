FROM golang:latest

RUN go get github.com/constabulary/gb/...

RUN mkdir -p /opt/oden
RUN curl -L https://github.com/oden-lang/oden/releases/download/0.3.0-alpha13/oden-0.3.0-alpha13-linux.tar.gz | tar -xvz -C /opt
RUN ln -s /opt/oden/bin/oden /usr/bin/oden

RUN mkdir -p /app
WORKDIR /app

COPY . /app
RUN gb build

EXPOSE 8080

ENTRYPOINT ["bin/playground"]
