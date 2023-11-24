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

COPY --from=BuildStage /myapp /myapp

CMD [ "/myapp" ]