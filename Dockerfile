FROM debian:bullseye-slim AS tailwind-builder

WORKDIR /app

# Install curl
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

# Download and verify Tailwind binary based on architecture
ARG TARGETARCH
RUN set -ex && \
    ARCH=$([ "$TARGETARCH" = "arm64" ] && echo "arm64" || echo "x64") && \
    TAILWIND_URL="https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-${ARCH}" && \
    curl -sL -o tailwindcss ${TAILWIND_URL} && \
    chmod +x tailwindcss && \
    ls -la tailwindcss

COPY static/css/input.css ./static/css/input.css

RUN ./tailwindcss -i ./static/css/input.css -o ./static/css/output.css

FROM golang:bullseye AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=tailwind-builder /app/static/css/output.css ./static/css/output.css

RUN CGO_ENABLED=0 GOOS=linux go build -o main

FROM gcr.io/distroless/base-debian11

WORKDIR /app

COPY --from=build-stage /app/main /main
COPY --from=build-stage /app/migrations /migrations
COPY --from=build-stage /app/static /static

EXPOSE 8080

ENTRYPOINT [ "/main" ]