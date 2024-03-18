FROM golang:1.21-alpine as BuildStage

# BUILD
WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /myapp

# RUN
FROM alpine

COPY .env /app/.env

WORKDIR /app

COPY cmd/migrate/migrations /app/cmd/migrate/migrations

COPY --from=BuildStage /myapp /myapp

CMD [ "/myapp" ]