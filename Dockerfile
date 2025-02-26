ARG TARGETOS
ARG TARGETARCH

FROM golang:1.23.6-alpine3.21 as base

# See https://stackoverflow.com/a/55757473/4752298
ENV USER=appuser
ENV UID=12345

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/does_not_exist" \
    --no-create-home \
    --shell "/sbin/nologin" \
    --uid "$UID" \
    "${USER}"

WORKDIR /src
COPY go.* ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go mod download
COPY .. .

FROM base AS build-webserver
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-w -s" -o /out/webserver cmd/webserver/main.go

FROM base as webserver
COPY --from=build-webserver /out/webserver /webserver
USER ${USER}:${USER}
ENTRYPOINT ["/webserver"]