# Team: Romanshk Volkov - https://github.com/RomanshkVolkov
# Team: Diegode - https://github.com/diegode-tsx
# Team: Alexandergv2117 - https://github.com/Alexandergv2117
# Customer: Dwit México - https://dwitmexico.com
# Project: <binary-name>

FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ARG DEPLOY_DB_DSN_BETA_AUTOPARTES
ENV DB_DSN_BETA_AUTOPARTES="$DEPLOY_DB_DSN_BETA_AUTOPARTES"

RUN go build -o /ws-beta-autopartes ./cmd/

ENV DEPLOY_DB_DSN_BETA_AUTOPARTES="hidden"

# runner light enviroment
FROM alpine:latest

WORKDIR /srv

COPY --from=builder /ws-beta-autopartes .
COPY ./static ./static

EXPOSE 8080

CMD ["./ws-beta-autopartes"]
