#RUN apt-get install git && apt-get mercurial && apt-get golang

FROM ubuntu:12.04
MAINTAINER Mallika Sen

# Mercurial
RUN echo 'deb http://ppa.launchpad.net/mercurial-ppa/releases/ubuntu precise main' > /etc/apt/sources.list.d/mercurial.list
RUN echo 'deb-src http://ppa.launchpad.net/mercurial-ppa/releases/ubuntu precise main' >> /etc/apt/sources.list.d/mercurial.list
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 323293EE

RUN apt-get update
RUN apt-get install -y curl git bzr mercurial

#install go
RUN curl -s https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz | tar -v -C /usr/local/ -xz

#Configure environment for Go
ENV PATH  /usr/local/go/bin:/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin
ENV GOPATH  /go
ENV GOROOT  /usr/local/go

#get project files from github
RUN go get github.com/AaronGoldman/ccfs
RUN apt-get -y install fuse 

#WORKDIR /go/src/github.com/AaronGoldman/ccfs
#ADD . /go/src/github.com/AaronGoldman/ccfs/ccfs
#ADD . ~/Desktop/create_objects.sh
WORKDIR /ccfs
ADD ccfs ccfs/ccfs
ADD ccfs.test ccfs/
ADD create_objects.sh ccfs/
#COPY ./ccfs/bin/. bin/

#any other dependencies?

#install application
#RUN go build 

#Do we need to expose port?
EXPOSE 8080:8080

ENTRYPOINT /ccfs/create_objects.sh
#./ccfs -mount

#sudo usermod -a -G fuse

