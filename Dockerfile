# Build Stage: Build bot using the alpine image, also install doppler in it
FROM golang:1.20 AS builder
RUN apt-get update && apt-get upgrade -y && apt-get install build-essential -y
WORKDIR /app
COPY . .
RUN GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o out/GigaUserbot -ldflags="-w -s" .

# Run Stage: Run bot using the bot and doppler binary copied from build stage
FROM golang:1.20
COPY --from=builder /app/out/GigaUserbot /app/GigaUserbot
COPY --from=builder /app/changelog.json /app/changelog.json
CMD ["/app/GigaUserbot"]
