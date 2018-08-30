FROM ubuntu:16.04

RUN apt-get update -y && \
    apt-get install fping -y

RUN apt-get install python3 -y && \
    apt-get install python3-pip -y && \
    pip3 install influxdb

ADD . /usr/local/testvm

WORKDIR /usr/local/testvm

#CMD ["sleep","9000000000"]
CMD ["bee run"]