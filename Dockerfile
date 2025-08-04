FROM alpine:3.22.0 AS base

FROM base AS builder

ARG TARGETOS
ARG TARGETARCH
ARG KSS_VERSION

RUN apk add --no-cache curl

RUN curl -fsSL https://github.com/moonlight8978/kubernetes-schema-store/releases/download/v${KSS_VERSION}/kubernetes-schema-store_${TARGETOS}_${TARGET_ARCH}.tar.gz | tar -xzf - -C /tmp \
  && chmod +x /tmp/kss

FROM base

COPY --from=builder /tmp/kss /usr/local/bin/kss

CMD ["kss"]
