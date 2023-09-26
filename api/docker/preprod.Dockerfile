###########################
## STEP 1 - Getting .env ##
###########################
FROM amazon/aws-cli AS deps

ARG AWS_ACCESS_KEY_ID
ARG AWS_SECRET_ACCESS_KEY
ARG AWS_DEFAULT_REGION
ARG S3_BUCKET_NAME
ARG CI_COMMIT_BRANCH

WORKDIR /tmp

RUN aws configure set aws_access_key_id $AWS_ACCESS_KEY_ID
RUN aws configure set aws_secret_access_key $AWS_SECRET_ACCESS_KEY
RUN aws configure set default.region $AWS_DEFAULT_REGION
RUN aws configure set output json

RUN aws s3 cp s3://$S3_BUCKET_NAME/base_env/coreapi/$CI_COMMIT_BRANCH/.env ./

RUN sed -i "s|eden-v2-production|pre-production|g" .env
RUN sed -i "s|10.2.13.135:20000|mysql-master.preprod.svc.cluster.local:3306|g" .env
RUN sed -i "s|10.2.13.141|mongodb.preprod.svc.cluster.local|g" .env
RUN sed -i "s|broker01.edenfarm.tech|kafka.preprod.svc.cluster.local|g" .env
RUN sed -i "s|10.2.13.172|redis-sentinel.preprod.svc.cluster.local|g" .env
RUN sed -i "s|notifapi|clone.notifapi|g" .env
RUN sed -i "s|printapi|clone.printapi|g" .env
RUN sed -i "s|coreapi|clone.coreapi|g" .env

######################
## STEP 2 - Builder ##
######################
FROM golang:1.16-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates gcc libc-dev gcompat protoc make openssl tzdata && update-ca-certificates

ARG GITLAB_CI_ACCESS
ARG GITLAB_CI_TOKEN

ENV USER=pengguna
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /app/api
COPY . .

RUN git config --global url."https://$GITLAB_CI_ACCESS:$GITLAB_CI_TOKEN@git.edenfarm.id/".insteadOf "https://git.edenfarm.id/"
RUN git clone -b $CI_COMMIT_BRANCH https://$GITLAB_CI_ACCESS:$GITLAB_CI_TOKEN@git.edenfarm.id/project-version2/datamodel.git ../datamodel
RUN go mod vendor

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./api-core

###########################
## STEP 3 - Distribution ##
###########################
FROM alpine:latest

RUN apk update && apk add bash

WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /app/api/api-core ./
COPY --from=deps /tmp/.env ./

USER pengguna:pengguna

CMD ["./api-core"]

