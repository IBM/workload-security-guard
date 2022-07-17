# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /protector cmd/protector/*.go


##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /protector /protector

EXPOSE 22000

USER nonroot:nonroot

ENTRYPOINT ["/protector"]
