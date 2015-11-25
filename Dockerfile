FROM heroku/go:1.5

WORKDIR /app/oden
RUN curl -L https://github.com/oden-lang/oden/releases/download/0.1.8/oden-0.1.8-linux.tar.gz | tar -xvz -C /app/oden

RUN mkdir -p /app/bin
RUN ln -s /app/oden/bin/odenc /app/user/bin/odenc
RUN ln -s /app/oden/bin/odenc /usr/bin/odenc
