FROM golang:1.20
COPY giga changelog.json /app/
CMD ["/app/giga"]
