#!/bin/bash
set -e -x
set -o pipefail

yum remove -y sendmail postfix
yum install -y wget

# INSERT ENV.SH HERE


####


# install
wget -q https://github.com/cfstras/textmail/releases/download/v1/textmail
chmod a+x textmail

cat > goguerrilla.conf <<EOF
{
  "GM_ALLOWED_HOSTS": "__MAIL_HOSTNAME__",
  "GM_MAX_CLIENTS": "500",
  "GM_PRIMARY_MAIL_HOST": "__MAIL_HOSTNAME__",
  "GSMTP_HOST_NAME": "__MAIL_HOSTNAME__",
  "GSMTP_LOG_FILE": "/goguerilla.log",
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
sed -i "s/__MAIL_HOSTNAME__/$MAIL_HOSTNAME/g" goguerrilla.conf


cat > /etc/logrotate.d/goguerilla <<EOF
/goguerilla.log {
    missingok
    copytruncate
    notifempty
    weekly
    rotate 10
    create 0600 root root
}
EOF

./textmail

