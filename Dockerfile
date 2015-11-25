FROM heroku/go:1.5

WORKDIR /app/oden
# Replace this with downloading a Linux build
RUN curl -L https://github.com/oden-lang/oden/releases/download/0.1.8/oden-0.1.8-linux.tar.gz | tar -xvz -C /app/oden

RUN ln -s /app/oden/bin/odenc /app/user/bin/odenc
