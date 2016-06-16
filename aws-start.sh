#!/bin/bash
set -e -x
set -o pipefail

yum remove -y sendmail
yum install -y curl git wget gcc python-pip jq
pip install awscli

export PATH=/opt/aws/bin:/opt/aws/apitools/ec2/bin:$PATH
export EC2_HOME=/opt/aws/apitools/ec2
export JAVA_HOME=/usr/lib/jvm/jre
export AWS_DEFAULT_REGION=us-west-2
export EC2_URL=ec2.$AWS_DEFAULT_REGION.amazonaws.com
export MAIL_HOSTNAME=txt.cfs.im

rm -rf go*.linux-amd64.tar.gz*
wget -q https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz
rm -rf /usr/local/go
tar -C /usr/local -xzf go1.6.2.linux-amd64.tar.gz
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin
export GOBIN=$GOROOT/bin
mkdir -p ~/go/ 
export GOPATH=~/go/
export PATH=$GOPATH/bin:$PATH 

# INSERT ENV.SH HERE

####

# install
go get github.com/cfstras/textmail
cd $GOPATH/src/github.com/cfstras/textmail
./start-app.sh

while true; do sleep 60;done

