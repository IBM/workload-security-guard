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
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /guard-learner /guard-learner

EXPOSE 80

USER nonroot:nonroot

ENTRYPOINT ["/guard-learner"]
