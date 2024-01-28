FROM golang:alpine AS build
RUN apk add git
ENV GOPATH=/go GO111MODULE=off CGO_ENABLED=0
ARG NEXCOMMON_TOKEN="ATCTT3xFfGN03BizwSyPDKLziZWaxfC7nCV56pdIcegNz7YBp8L3zqY0FqrwM9c3UDVL2wXlsxvgUUHFYK20zUJomLa3ndZvWDre0_oLqpmWWXgEeBgOEkjr2Q4vgLfIKok-m8XQvhr6sK6oPXRKfRGB6J7rogDTFVEeqyV_VWK9Yekv1B8EBoM=5E6A208E"
ARG NEXCOMMON_BRANCH

WORKDIR /go/src/nexsoft.co.id
RUN git clone https://x-token-auth:${NEXCOMMON_TOKEN}@bitbucket.org/nexcloud/nexcommon.git -b ${NEXCOMMON_BRANCH} --single-branch
COPY . nextrac2
RUN go get ./nextrac2/...
RUN go build -ldflags '-s' -o /home/nextrac2 ./nextrac2/main.go

FROM alpine
RUN apk add curl
RUN apk add tz

ENV TZ="Asia/Jakarta"

RUN mkdir -p /home/config
RUN mkdir -p /home/apptrac/go/nextrac/i18n
RUN mkdir -p /home/apptrac/go/nextrac/logs
RUN mkdir -p /home/sql_migrations
RUN mkdir -p /home/v1
RUN mkdir -p /home/apptrac/go/nextrac/data
RUN mkdir -p /home/data
RUN mkdir -p /home/apptrac/go/nextrac/cdn/nextrac

WORKDIR /home
COPY --from=0 /home/nextrac2 ./
COPY i18n /home/apptrac/go/nextrac/i18n
COPY config config
COPY sql_migrations sql_migrations

EXPOSE 9000

# LOGGING
RUN touch /home/apptrac/go/nextrac/logs/nextrac-stdout.log
RUN ln -sf /proc/1/fd/1 /home/apptrac/go/nextrac/logs/nextrac-stdout.log

# PORTAINER
LABEL com.centurylinklabs.watchtower.enable="true"
LABEL io.portainer.accesscontrol.teams="nextrac"
LABEL io.portainer.accesscontrol.users="tinus"
HEALTHCHECK --interval=1m30s --timeout=30s --start-period=30s --retries=3 CMD [ "curl", "http://localhost:9000/v1/nextrac/health" ]