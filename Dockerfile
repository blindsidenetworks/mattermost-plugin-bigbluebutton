FROM golang:1.13-alpine

RUN apk add --no-cache bash git

RUN apk add --virtual build-dependencies \
        build-base \
        gcc \
        wget \
        curl

ENV VERSION=v14.3.0 NPM_VERSION=6 YARN_VERSION=latest

RUN apk upgrade --no-cache -U && \
  apk add --no-cache gnupg libstdc++

RUN curl -sfSLO https://unofficial-builds.nodejs.org/download/release/${VERSION}/node-${VERSION}-linux-x64-musl.tar.xz && \
  curl -sfSLO https://unofficial-builds.nodejs.org/download/release/${VERSION}/SHASUMS256.txt && \
  grep " node-${VERSION}-linux-x64-musl.tar.xz\$" SHASUMS256.txt | sha256sum -c | grep ': OK$' && \
  tar -xf node-${VERSION}-linux-x64-musl.tar.xz -C /usr --strip 1 && \
  rm node-${VERSION}-linux-x64-musl.tar.xz

RUN npm install -g npm@${NPM_VERSION} && \
  find /usr/lib/node_modules/npm -type d \( -name test -o -name .bin \) | xargs rm -rf

RUN for server in ipv4.pool.sks-keyservers.net keyserver.pgp.com ha.pool.sks-keyservers.net; do \
    gpg --keyserver $server --recv-keys \
      6A010C5166006599AA17F08146C2130DFD2497F5 && break; \
  done && \
  curl -sfSL -O https://yarnpkg.com/${YARN_VERSION}.tar.gz -O https://yarnpkg.com/${YARN_VERSION}.tar.gz.asc && \
  gpg --batch --verify ${YARN_VERSION}.tar.gz.asc ${YARN_VERSION}.tar.gz && \
  mkdir /usr/local/share/yarn && \
  tar -xf ${YARN_VERSION}.tar.gz -C /usr/local/share/yarn --strip 1 && \
  ln -s /usr/local/share/yarn/bin/yarn /usr/local/bin/ && \
  ln -s /usr/local/share/yarn/bin/yarnpkg /usr/local/bin/ && \
  rm ${YARN_VERSION}.tar.gz*

ENV PYTHONUNBUFFERED=1

RUN apk add --no-cache python && \
    python -m ensurepip && \
    rm -r /usr/lib/python*/ensurepip && \
    pip install --upgrade pip setuptools && \
    rm -r /root/.cache

RUN apk add --no-cache expat expat-dev

#COPY /home/basti/daten/Projekte/bigbluebutton/workspace/go-src-cache /go/src/

RUN mkdir -p /go/src/github.com/blindsidenetworks/mattermost-plugin-bigbluebutton

WORKDIR /go/src/github.com/blindsidenetworks/mattermost-plugin-bigbluebutton

COPY . /go/src/github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/

RUN make build
