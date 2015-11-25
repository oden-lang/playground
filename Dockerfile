FROM heroku/go:1.5

WORKDIR /tmp/
RUN curl -L https://github.com/oden-lang/oden/releases/download/0.1.7/oden-0.1.7-osx.tar.gz > oden.tar.gz
