FROM golang:bookworm
LABEL authors="Mugagga k"
LABEL org="H26"
LABEL version="1.0.0"
WORKDIR /app
COPY . .
RUN go mod download
EXPOSE 2556:2556
CMD [ "go", "run", "main.go" ]
