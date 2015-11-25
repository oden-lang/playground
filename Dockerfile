FROM heroku/go:1.5

WORKDIR /tmp/

# Replace this with downloading a Linux build
RUN curl -L https://github.com/oden-lang/oden/releases/download/0.1.7/oden-0.1.7-osx.tar.gz > oden.tar.gz

# todo: Unpack to /app/oden

RUN echo "export PATH=\"/app/oden/bin:\$PATH\"" >> /app/.profile.d/gm.sh
