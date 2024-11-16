FROM golang:1.21.6
WORKDIR /gnats
COPY . .
CMD ["go", "run", "./main.go"]
