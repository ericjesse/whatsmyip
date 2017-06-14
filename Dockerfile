FROM alpine:latest

MAINTAINER Eric Jessé <eric@jesse.fr>

WORKDIR "/opt"

ADD .docker_build/whatsmyip /opt/bin/whatsmyip

ENV PORT 5000
EXPOSE 5000

CMD ["/opt/bin/whatsmyip"]
