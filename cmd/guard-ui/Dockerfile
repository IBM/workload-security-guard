# syntax=docker/dockerfile:1

##
## Frontend Build
##
FROM node:latest as fe 
WORKDIR /

COPY cmd/guard-ui/frontend/ ./
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
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /guard-ui /guard-ui
COPY --from=fe /build /frontend/build

EXPOSE 80

USER nonroot:nonroot

ENTRYPOINT ["/guard-ui"]
