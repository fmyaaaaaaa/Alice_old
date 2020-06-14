# build
FROM golang:1.14-alpine3.12 as build

RUN apk add --update --no-cache git
WORKDIR /go/src/github.com/fmyaaaaaaa/Alice/alice-trading
COPY . .
RUN go build -o alice .

# copy
FROM alpine:3.12.0

RUN apk add --update --no-cache ca-certificates
WORKDIR /go/src/github.com/fmyaaaaaaa/Alice/alice-trading
COPY --from=build /go/src/github.com/fmyaaaaaaa/Alice/alice-trading/alice /go/src/github.com/fmyaaaaaaa/Alice/alice-trading/alice

# 環境変数設定
ARG MODE
ARG URL
ARG ACCOUNT_ID
ARG ACCESS_TOKEN
ARG DB_HOST
ARG DB_PORT
ARG DB_USERNAME
ARG DB_PASSWORD
ARG DB_NAME
ENV MODE=$MODE URL=$URL ACCOUNT_ID=$ACCOUNT_ID ACCESS_TOKEN=$ACCESS_TOKEN DB_HOST=$DB_HOST DB_PORT=$DB_PORT DB_USERNAME=$DB_USERNAME DB_PASSWORD=$DB_PASSWORD DB_NAME=$DB_NAME
ENTRYPOINT /go/src/github.com/fmyaaaaaaa/Alice/alice-trading/alice $MODE $URL $ACCOUNT_ID $ACCESS_TOKEN $DB_HOST $DB_PORT $DB_USERNAME $DB_PASSWORD $DB_NAME
