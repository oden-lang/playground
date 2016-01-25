FROM heroku/go-gb:1.5.1

RUN mkdir -p /app/bin
RUN mkdir -p /app/user/bin

RUN mkdir -p /app/oden
RUN curl -L https://github.com/oden-lang/oden/releases/download/0.2.1-RC2/oden-0.2.1-RC2-linux.tar.gz | tar -xvz -C /app/oden
RUN ln -s /app/oden/bin/odenc /app/bin/odenc
RUN ln -s /app/oden/bin/odenc /app/user/bin/odenc


RUN mkdir -p /app/user \
      && cp -r /app/.cache/go /app/user/go
WORKDIR /app
