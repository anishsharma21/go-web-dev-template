FROM golang:bullseye AS build-stage

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v4.0.7/tailwindcss-linux-x64 && \
    chmod +x tailwindcss-linux-x64 && \
    mv tailwindcss-linux-x64 tailwindcss

RUN ./tailwindcss -i ./static/css/input.css -o ./static/css/output.css

RUN CGO_ENABLED=0 GOOS=linux go build -o main

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build-stage /main /main
COPY --from=build-stage /migrations /migrations
COPY --from=build-stage /static /static

EXPOSE 8080

ENTRYPOINT [ "/main" ]