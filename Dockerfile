# Build Stage: Build bot using the alpine image, also install doppler in it
FROM golang AS builder
RUN apt-get update && apt-get upgrade -y && apt-get install build-essential -y
WORKDIR /app
COPY . .
RUN GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o out/GigaUserbot -ldflags="-w -s" .

# Run Stage: Run bot using the bot and doppler binary copied from build stage
FROM ubuntu
COPY --from=builder /app/out/GigaUserbot /app/GigaUserbot
COPY config.json config.json
CMD ["/app/GigaUserbot"]
