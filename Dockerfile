ARG BUILD_VERSION="0.4.1"

FROM golang:1.22-alpine3.19 as builder
ARG BUILD_VERSION

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN touch config.txt
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstraper -ldflags="-X 'github.com/tigerwill90/djdiscord/internal/build.Version=${BUILD_VERSION}'" main.go

FROM gcr.io/distroless/java17-debian12
ARG BUILD_VERSION

COPY JMusicBot-${BUILD_VERSION}.jar .
COPY --from=builder --chown=65532:65532 /app/bootstraper .
COPY --from=builder --chown=65532:65532 /app/Playlists /Playlists
COPY --from=builder --chown=65532:65532 /app/config.txt .

ENTRYPOINT ["/bootstraper"]

USER 65532
