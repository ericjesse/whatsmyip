FROM alpine:latest

MAINTAINER Eric Jessé <eric@jesse.fr>

WORKDIR "/opt"

ADD .docker_build/whatsmyip /opt/bin/whatsmyip

CMD ["/opt/bin/whatsmyip"]

