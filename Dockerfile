# Vaku Rust runtime image. Expects a statically-linked vaku binary in the
# build context (the release workflow drops it next to this Dockerfile).
FROM scratch
ARG TARGETPLATFORM

LABEL org.opencontainers.image.ref.name="vaku" \
    org.opencontainers.image.title="vaku" \
    org.opencontainers.image.description="A CLI to extend the official Vault client" \
    org.opencontainers.image.licenses="MIT" \
    org.opencontainers.image.authors="sean@lingren.com" \
    org.opencontainers.image.url="https://vaku.dev" \
    org.opencontainers.image.documentation="https://vaku.dev" \
    org.opencontainers.image.source="https://github.com/lingrino/vaku"

COPY $TARGETPLATFORM/vaku /
ENTRYPOINT ["/vaku"]
