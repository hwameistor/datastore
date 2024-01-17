FROM centos:7

RUN yum upgrade nss -y

COPY ./_build/metadata-controller /

ENTRYPOINT [ "/metadata-controller" ]