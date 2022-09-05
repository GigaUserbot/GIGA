FROM golang:1.19
COPY giga changelog.json /app/
CMD ["/app/giga"]
