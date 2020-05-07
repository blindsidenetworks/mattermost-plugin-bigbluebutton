FROM golang:1.14

RUN apt update && \
    apt -y install build-essential npm

RUN addgroup --gid 1000 node \
    && useradd --create-home --uid 1000 --gid node --shell /bin/sh node

RUN npm update && npm cache clean -f && npm install -g n && n 10

CMD /bin/sh
