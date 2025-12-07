FROM golang:1.24-alpine AS builder

WORKDIR /build

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

ADD go.mod .
ADD go.sum .

RUN go mod download

COPY . .

RUN go build -o /app/todo-app ./main.go

FROM alpine:latest

RUN addgroup -S app && adduser -S todo-app -G app
RUN mkdir -p /app/web /app/data
RUN chown -R todo-app:app /app

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/data/scheduler.db
ENV TODO_PASSWORD=""

WORKDIR /app

COPY --from=builder /app/todo-app /app/todo-app
COPY web /app/web

RUN chmod +x /app/todo-app
RUN chown -R todo-app:app /app

USER todo-app

EXPOSE 7540

ENTRYPOINT ["/app/todo-app"]
