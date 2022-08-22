# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /guard-learner cmd/guard-learner/*.go



##
## Deploy
##
FROM gcr.io/distroless/base-alpine

WORKDIR /

# Install the ibmcloud cli
RUN curl -fsSL https://clis.cloud.ibm.com/install/linux | sh && \
    ibmcloud plugin install container-service -f && \
    ibmcloud plugin install container-registry -f && \
    ibmcloud plugin install observe-service -f && \
    ibmcloud plugin install code-engine -f && \
    if [[ "$TARGETPLATFORM" != "linux/arm64" ]]; then ibmcloud plugin install vpc-infrastructure -f; fi && \
    ibmcloud config --check-version=false && \
    chmod -R g=u ${HOME}

COPY --from=build /guard-learner /guard-learner

EXPOSE 8888

USER nonroot:nonroot

ENTRYPOINT ["/guard-learner"]
