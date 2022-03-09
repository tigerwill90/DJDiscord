FROM golang:1.17.8-alpine3.15 as builder
ARG BUILD_VERSION="0.3.6"

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN touch config.txt
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bot -ldflags="-X 'djdiscord/internal/build.Version=${BUILD_VERSION}'" main.go

FROM gcr.io/distroless/java17-debian11

COPY JMusicBot-0.3.6.jar .
COPY --from=builder /app/bot .
COPY --from=builder --chown=65532:65532 /app/Playlists /Playlists
COPY --from=builder --chown=65532:65532 /app/config.txt .

ENTRYPOINT ["/bot"]

USER 65532