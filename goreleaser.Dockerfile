FROM golang:1.19
COPY GIGA changelog.json /app/
CMD ["/app/GigaUserbot"]
