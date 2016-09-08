FROM ubuntu
RUN apt-get update && apt-get -y install python-software-properties 
RUN apt-get update && apt-get -y install apt-file && apt-file update && apt-file search add-apt-repository && apt-get -y install software-properties-common
RUN apt-get -y install libcurl4-gnutls-dev libexpat1-dev gettext libz-dev libssl-dev git
RUN yes | add-apt-repository ppa:ubuntu-lxc/lxd-stable && apt-get update && apt-get -y install golang
ENV GOPATH /go 
ENV PATH=$PATH:/go/bin
RUN go get github.com/tools/godep && go get github.com/boiledgas/receiver-test
RUN cd /go/src/github.com/boiledgas/receiver-test && godep restore
RUN cd /go/src/github.com/boiledgas/receiver-test && go install github.com/boiledgas/receiver-test/main

CMD ["/go/bin/main"]