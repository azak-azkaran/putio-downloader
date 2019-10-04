FROM golang:alpine AS build-env
RUN apk --no-cache add git make
RUN go get  github.com/azak-azkaran/putio-downloader
WORKDIR /go/src/github.com/azak-azkaran/putio-downloader
RUN make install

FROM alpine
# less priviledge user, the id should map the user the downloaded files belongs to
RUN apk --no-cache add shadow && \
        groupadd -r dummy && \
        useradd -r -g dummy dummy -u 1000

WORKDIR /opt/putio-downloader/

COPY --from=build-env /go/bin/putio-downloader /opt/putio-downloader/main

ENTRYPOINT ["/opt/putio-downloader/main"]
