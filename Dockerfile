FROM heroku/go:1.5

RUN apt-get install golang

WORKDIR /app/oden
RUN curl -L https://github.com/oden-lang/oden/releases/download/0.1.9/oden-0.1.9-linux.tar.gz | tar -xvz -C /app/oden

RUN mkdir -p /app/bin
RUN ln -s /app/oden/bin/odenc /app/bin/odenc
RUN mkdir -p /app/user/bin
RUN ln -s /app/oden/bin/odenc /app/user/bin/odenc
