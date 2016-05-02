#!/bin/bash
set -e -x
set -o pipefail

cat > goguerrilla.conf <<EOF
{
  "GM_ALLOWED_HOSTS": "__MAIL_HOSTNAME__,__HOSTNAME__",
  "GM_MAX_CLIENTS": "500",
  "GM_PRIMARY_MAIL_HOST": "__MAIL_HOSTNAME__",
  "GSMTP_HOST_NAME": "__HOSTNAME__",
  "GSMTP_LOG_FILE": "",
  "GSMTP_MAX_SIZE": "131072",
  "GSMTP_PRV_KEY": "",
  "GSMTP_PUB_KEY": "",
  "GSMTP_TIMEOUT": "100",
  "GSMTP_VERBOSE": "Y",
  "GSTMP_LISTEN_INTERFACE": "0.0.0.0:25",
  "NGINX_AUTH": "127.0.0.1:8025",
  "NGINX_AUTH_ENABLED": "N",
  "SGID": "1000",
  "SUID": "1000"
}
EOF

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

