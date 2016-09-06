#!/bin/bash
set -e -x
set -o pipefail



REGION=`curl http://169.254.169.254/latest/dynamic/instance-identity/document|grep region|awk -F\" '{print $4}'`

HOSTNAME=$(curl -s http://169.254.169.254/latest/meta-data/public-hostname)
sed -i "s/__HOSTNAME__/$HOSTNAME/g" goguerrilla.conf
sed -i "s/__MAIL_HOSTNAME__/$MAIL_HOSTNAME/g" goguerrilla.conf

ALLOC=$(ec2-describe-addresses --region $REGION --show-empty-fields $(dig +short txt.cfs.im) | awk '{print $5}')

ec2-associate-address \
  -i $(curl -s http://169.254.169.254/latest/meta-data/instance-id) \
  -a $ALLOC \
  --allow-reassociation \
  --region $REGION

echo "Okay! send mails to $MAIL_AUTH_PREFIX-<number>@$MAIL_HOSTNAME !"

textmail

