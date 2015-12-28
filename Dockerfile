FROM heroku/go-gb:1.5.1

RUN mkdir -p /app/bin
RUN mkdir -p /app/user/bin

RUN mkdir -p /app/oden
RUN curl -L https://github.com/oden-lang/oden/releases/download/0.2.0/oden-0.2.0-linux.tar.gz | tar -xvz -C /app/oden
RUN ln -s /app/oden/bin/odenc /app/bin/odenc
RUN ln -s /app/oden/bin/odenc /app/user/bin/odenc

WORKDIR /app
RUN curl -L https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz | tar -xvz go/bin/gofmt

RUN cp /app/go/bin/gofmt /app/bin/gofmt
RUN cp /app/go/bin/gofmt /app/user/bin/gofmt
RUN rm -r /app/go
