FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/ratelimit-operator

ADD build/_output/bin/ratelimit-operator ${OPERATOR}

CMD ${OPERATOR}
