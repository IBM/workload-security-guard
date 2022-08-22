# syntax=docker/dockerfile:1

##
## Frontend Build
##
FROM node:latest as fe 
WORKDIR /

COPY cmd/guard-ui/frontend/ ./
RUN npm install
RUN npm run build

##
## Build
##
FROM golang:1.18 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /guard-ui cmd/guard-ui/*.go

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

COPY --from=build /guard-ui /guard-ui
COPY --from=fe /build /frontend/build

EXPOSE 80

USER nonroot:nonroot
