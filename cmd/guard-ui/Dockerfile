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
FROM alpine/curl

WORKDIR /

# Install the ibmcloud cli
RUN curl -fsSL https://clis.cloud.ibm.com/install/linux | sh
RUN ibmcloud update
RUN ibmcloud plugin install container-service -f
RUN ibmcloud plugin install container-registry -f
RUN ibmcloud plugin install observe-service -f
RUN ibmcloud plugin install code-engine -f
RUN ibmcloud config --check-version=false

COPY --from=build /guard-ui /guard-ui
COPY --from=fe /build /frontend/build

EXPOSE 80

RUN adduser -D nonroot
USER nonroot:nonroot
