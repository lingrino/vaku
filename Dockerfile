FROM scratch
LABEL maintainer="sean@lingrino.com"
COPY vaku /
ENTRYPOINT ["/vaku"]
