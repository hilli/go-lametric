FROM golang:1.22.4-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app
COPY . /app
RUN go mod download
WORKDIR /app/cmd/lamertric-homekit-hub
RUN go mod download
RUN go build -o /app/lametric-homekit-hub .

FROM alpine:latest as production
WORKDIR /app
COPY --from=builder /app/lametric-homekit-hub /app/lametric-homekit-hub

RUN apk --no-cache add curl

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/app/lametric-homekit-hub"]

# use curl to do a health check
HEALTHCHECK --interval=5s --timeout=3s --retries=3 CMD curl -f http://localhost:8080/_ping || exit 1