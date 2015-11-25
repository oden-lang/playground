FROM heroku/go:1.5

WORKDIR /tmp/oden

# Replace this with downloading a Linux build
RUN curl -L https://github.com/oden-lang/oden/releases/download/0.1.8/oden-0.1.8-linux.tar.gz | tar -xvz -C /tmp/oden

RUN echo "export PATH=\"/app/oden/bin:\$PATH\"" >> /app/.profile.d/gm.sh
