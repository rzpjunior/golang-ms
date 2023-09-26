FROM git.edenfarm.id:5050/image/golang-edenfarm/master:latest as build

ARG AWS_ACCESS_KEY_ID
ARG AWS_SECRET_ACCESS_KEY
ARG GITLAB_CI_ACCESS
ARG GITLAB_CI_TOKEN
ARG S3_BUCKET_NAME
ARG AWS_DEFAULT_REGION
ARG CI_COMMIT_BRANCH
ARG CI_COMMIT_TAG

RUN mkdir /go/src/git.edenfarm.id
RUN mkdir /go/src/git.edenfarm.id/project-version2
RUN git config --global url."https://$GITLAB_CI_ACCESS:$GITLAB_CI_TOKEN@git.edenfarm.id/".insteadOf "https://git.edenfarm.id/"
RUN git clone https://$GITLAB_CI_ACCESS:$GITLAB_CI_TOKEN@git.edenfarm.id/project-version2/api.git /go/src/git.edenfarm.id/project-version2/api
RUN git clone https://$GITLAB_CI_ACCESS:$GITLAB_CI_TOKEN@git.edenfarm.id/project-version2/datamodel.git /go/src/git.edenfarm.id/project-version2/datamodel
RUN cd /go/src/git.edenfarm.id/project-version2/datamodel && git checkout $CI_DEFAULT_BRANCH && git add . && git stash && git pull origin $CI_DEFAULT_BRANCH

WORKDIR /go/src/git.edenfarm.id/project-version2/api

RUN git checkout $CI_DEFAULT_BRANCH && git add . && git stash && git pull origin $CI_DEFAULT_BRANCH
RUN go mod vendor
RUN aws configure set aws_access_key_id $AWS_ACCESS_KEY_ID
RUN aws configure set aws_secret_access_key $AWS_SECRET_ACCESS_KEY
RUN aws configure set default.region $AWS_DEFAULT_REGION
RUN aws configure set output json

RUN aws s3 cp s3://$S3_BUCKET_NAME/base_env/coreapi/$CI_DEFAULT_BRANCH/.env ./
RUN sed -i "s/<VERSION>/${CI_COMMIT_TAG}/g" .env
RUN go build -o api .

FROM git.edenfarm.id:5050/image/golang-edenfarm/master:latest as deploy

WORKDIR /go/src/git.edenfarm.id/project-version2/api

COPY --from=build /go/src/git.edenfarm.id/project-version2/api/api /go/src/git.edenfarm.id/project-version2/api/api
COPY --from=build /go/src/git.edenfarm.id/project-version2/api/.env /go/src/git.edenfarm.id/project-version2/api/.env

CMD ["./api"]