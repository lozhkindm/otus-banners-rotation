# Собираем в гошке
FROM golang:1.17 as build

ENV BIN_FILE /opt/banners/goose
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/goose/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="goose"
LABEL MAINTAINERS="lozhkindm@yandex.ru"

ENV BIN_FILE "/opt/banners/goose"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ARG CONFIG_FILE_NAME
ENV CONFIG_FILE /etc/banners/.env
COPY ./configs/${CONFIG_FILE_NAME} ${CONFIG_FILE}

ARG DIR
COPY ./${DIR} /${DIR}

CMD ${BIN_FILE} --dir=/${DIR} --config=${CONFIG_FILE} up
