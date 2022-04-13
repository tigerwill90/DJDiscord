ARG BUILD_VERSION="0.3.8"

FROM golang:1.18.1-alpine3.15 as builder
ARG BUILD_VERSION

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN touch config.txt
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap -ldflags="-X 'djdiscord/internal/build.Version=${BUILD_VERSION}'" main.go

FROM gcr.io/distroless/java17-debian11
ARG BUILD_VERSION

COPY JMusicBot-${BUILD_VERSION}.jar .
COPY --from=builder --chown=65532:65532 /app/bootstrap .
COPY --from=builder --chown=65532:65532 /app/Playlists /Playlists
COPY --from=builder --chown=65532:65532 /app/config.txt .

ENTRYPOINT ["/bootstrap"]

USER 65532