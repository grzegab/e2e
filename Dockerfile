# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS base

COPY ./e2e /go/src/e2e
WORKDIR /go/src/e2e

ENV PATH="$PATH:/go/src/e2e"

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/src/e2e/e2e-tests

