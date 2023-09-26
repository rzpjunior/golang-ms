FROM git.edenfarm.id:5050/image/golang-18-alpine/master:18-5ca529ba AS builder

ARG CI_COMMIT_BRANCH
ARG GITLAB_CI_ACCESS
ARG GITLAB_CI_TOKEN

WORKDIR /app/api
COPY . .

RUN git clone -b $CI_COMMIT_BRANCH https://$GITLAB_CI_ACCESS:$GITLAB_CI_TOKEN@git.edenfarm.id/edenlabs/edenlabs.git ../edenlabs
RUN git clone -b $CI_COMMIT_BRANCH https://$GITLAB_CI_ACCESS:$GITLAB_CI_TOKEN@git.edenfarm.id/project-version3/erp-pkg/erp-client-grpc.git ../erp-client-grpc
RUN git clone -b $CI_COMMIT_BRANCH https://$GITLAB_CI_ACCESS:$GITLAB_CI_TOKEN@git.edenfarm.id/project-version3/erp-services/erp-protobuf.git ../erp-protobuf

RUN go mod tidy
RUN go mod vendor

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./service

FROM alpine

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/api/service ./
