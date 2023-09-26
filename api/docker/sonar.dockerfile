FROM $CI_REGISTRY/image/golang-edenfarm/master:latest as build

ARG CI_REGISTRY
ARG AWS_ACCESS_KEY_ID
ARG AWS_SECRET_ACCESS_KEY
ARG GITLAB_CI_ACCESS
ARG GITLAB_CI_TOKEN
ARG S3_BUCKET_NAME
ARG AWS_DEFAULT_REGION
ARG CI_COMMIT_BRANCH
ARG REPO_DIR
ARG CI_REPOSITORY_URL
ARG BASE_ENV_COREAPI
ARG CI_PROJECT_NAME
ARG SONARQUBE_URL
ARG SONARQUBE_PROJECT_KEY
ARG SONARQUBE_TOKEN

RUN git config --global url."https://$GITLAB_CI_ACCESS:$GITLAB_CI_TOKEN@git.edenfarm.id/".insteadOf "https://git.edenfarm.id/"

RUN git clone --single-branch --branch $CI_COMMIT_BRANCH $CI_REPOSITORY_URL $REPO_DIR/api
RUN git clone --single-branch --branch $CI_COMMIT_BRANCH https://$GITLAB_CI_ACCESS:$GITLAB_CI_TOKEN@git.edenfarm.id/project-version2/datamodel.git $REPO_DIR/datamodel

WORKDIR $REPO_DIR/api

RUN aws configure set aws_access_key_id $AWS_ACCESS_KEY_ID
RUN aws configure set aws_secret_access_key $AWS_SECRET_ACCESS_KEY
RUN aws configure set default.region $AWS_DEFAULT_REGION
RUN aws configure set output json

RUN aws s3 cp s3://$S3_BUCKET_NAME/$BASE_ENV_COREAPI/$CI_COMMIT_BRANCH/.env ./
RUN touch coverage.out
RUN go test ./... -coverprofile=coverage.out

FROM $CI_REGISTRY/image/sonar-scanner-cli/master:latest as upload

WORKDIR $REPO_DIR/api

COPY --from=build /go/src/git.edenfarm.id/project-version2/api/coverage.out /go/src/git.edenfarm.id/project-version2/api/coverage.out

RUN sonar-scanner -Dsonar.host.url=$SONARQUBE_URL -Dsonar.projectKey=$SONARQUBE_PROJECT_KEY -Dsonar.login=$SONARQUBE_TOKEN -Dsonar.projectName=Project-Core-$CI_PROJECT_NAME -Dsonar.sources=./src -Dsonar.go.coverage.reportPaths=coverage.out