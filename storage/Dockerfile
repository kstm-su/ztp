FROM ubuntu:14.04
RUN apt-get update
RUN apt-get install -y \
		apt-transport-https \
		ca-certificates \
		curl \
		software-properties-common \
		git
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
RUN apt-key fingerprint 0EBFCD88
RUN add-apt-repository \
		"deb [arch=amd64] https://download.docker.com/linux/ubuntu \
		$(lsb_release -cs) \
		stable"
RUN apt-get update
RUN apt-get install -y docker-ce

RUN curl -O https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz
RUN tar -xvf go1.8.linux-amd64.tar.gz
RUN mv go /usr/local
ENV PATH $PATH:/usr/local/go/bin
RUN go version
RUN echo 'export PATH=$PATH:/usr/local/go/bin' >> /root/.bash_profile

RUN go get -u github.com/moby/tool/cmd/moby
RUN go get -u github.com/julienschmidt/httprouter

RUN cd /root/go/src/github.com/moby/tool/cmd/moby \
	&& go build 
ENV PATH $PATH:/root/go/src/github.com/moby/tool/cmd/moby/
RUN echo 'export PATH=$PATH:/root/go/src/github.com/moby/tool/cmd/moby/' >> /root/.bash_profile
