FROM golang:1.24.6-bookworm AS build-stage
LABEL authors="Mugagga K"
LABEL org="H26"
LABEL version="1.0.0"
LABEL description="Authentication Service for Slim Redesigned"

# set the working directory and copy the go.mod and go.sum files
WORKDIR /app
COPY go.mod go.sum ./

# download the dependencies
RUN go mod download

#copy the application code
COPY . .

#build the app binary
RUN CGO_ENABLED=0 go build -o auth_service

FROM debian:bookworm-slim AS production-stage
WORKDIR /app
EXPOSE 2556

# copy certificates
COPY --from=build-stage /app/Certificates ./Certificates
COPY --from=build-stage /app/auth_service .

ENTRYPOINT [ "./auth_service" ]
