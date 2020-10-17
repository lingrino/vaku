FROM scratch

# https://github.com/opencontainers/image-spec/blob/master/annotations.md
LABEL org.opencontainers.image.ref.name="vaku" \
    org.opencontainers.image.ref.title="vaku" \
    org.opencontainers.image.description="A CLI to extend the official Vault client" \
    org.opencontainers.image.licenses="MIT" \
    org.opencontainers.image.authors="sean@lingrino.com" \
    org.opencontainers.image.url="https://vaku.dev" \
    org.opencontainers.image.documentation="https://vaku.dev" \
    org.opencontainers.image.source="https://github.com/lingrino/vaku"

COPY vaku /
ENTRYPOINT ["/vaku"]
