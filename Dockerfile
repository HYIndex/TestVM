FROM alpine

COPY . /usr/local/testvm

WORKDIR /usr/local/testvm


#CMD ['sh', 'build.sh']
#CMD ['/usr/local/testvm/bin/testvm']



CMD ["sleep","9000000000"]
