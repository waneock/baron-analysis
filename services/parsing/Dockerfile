FROM golang:1.26 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./services/parsing/cmd

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /server /server
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]