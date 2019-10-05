FROM golang:alpine AS build-env
# less priviledge user, the id should map the user the downloaded files belongs to
RUN apk --no-cache add git make shadow && \
        groupadd -r dummy && \
        useradd -r -g dummy dummy -u 1000
RUN go get  github.com/azak-azkaran/putio-downloader
WORKDIR /go/src/github.com/azak-azkaran/putio-downloader
COPY *.go ./
RUN make install

ENTRYPOINT ["/go/bin/putio-downloader"]
