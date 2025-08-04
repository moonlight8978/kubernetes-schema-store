FROM alpine:3.22.0 AS builder

ARG TARGETOS
ARG TARGETARCH
ARG KSS_VERSION

RUN apk add --no-cache curl

RUN curl -fsSL https://github.com/moonlight8978/kubernetes-schema-store/archive/refs/tags/v0.0.2.zip -o kss.tar.gz
