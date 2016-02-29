FROM heroku/go-gb:1.6

RUN mkdir -p /app/bin
RUN mkdir -p /app/user/bin

RUN mkdir -p /app/oden
RUN curl -L https://github.com/oden-lang/oden/releases/download/0.3.0-alpha11/oden-0.3.0-alpha11-linux.tar.gz | tar -xvz -C /app/oden
RUN ln -s /app/oden/bin/oden /app/bin/oden
RUN ln -s /app/oden/bin/oden /app/user/bin/oden

RUN mkdir -p /app/user \
      && cp -r /app/.cache/go /app/user/go
WORKDIR /app
