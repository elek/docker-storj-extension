FROM node:17.7-alpine3.14 AS client-builder
WORKDIR /app/web
ARG TARGETARCH
COPY web /app/web
RUN npm install
RUN npx vite build

FROM golang:1.18-alpine as backend-builder
WORKDIR /app
ADD . /app
RUN go build -o ./backend

FROM ghcr.io/elek/distribution:618d19fb as registry

FROM alpine
RUN apk add docker
LABEL org.opencontainers.image.title="Storj Decentralized Docker Registry" \
    org.opencontainers.image.description="An extension to start a local registry backed by decentralized Storj." \
    org.opencontainers.image.vendor="Storj Labs" \
    com.docker.desktop.extension.api.version=">= 0.2.0" \
    com.docker.desktop.extension.icon="https://assets.website-files.com/602eda09fc78afc76e9706b6/609177f5057ffc468b6ec24a_logo-mark.svg" \
    com.docker.extension.detailed-description="<h1>Description</h1><p>This extension helps to start a local docker registry which is backed by Storj decentralized storage.</p>" \
    com.docker.extension.publisher-url="https://storj.io" \
    com.docker.extension.screenshots='[{"alt":"main screenshot","url":"https://assets.website-files.com/602eda09fc78afc76e9706b6/609177f5057ffc468b6ec24a_logo-mark.svg"}]' \
    com.docker.extension.additional-urls='[{"title":"Storj decentralized cloud","url":"https://storj.io"}]' \
    com.docker.extension.changelog="<ul><li>Initial version</li></ul>"

RUN mkdir -p /run/guest-services/
COPY --from=client-builder /app/web/dist /ui
COPY --from=backend-builder /app/backend /backend
COPY --from=registry /bin/registry /usr/bin/registry
COPY config.yml .
COPY metadata.json .
ADD docker-compose.yaml .
ADD storj.svg .
CMD ["/backend/docker-storj-extension","run"]
