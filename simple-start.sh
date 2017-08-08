#!/bin/bash
set -ex
trap 'echo "exit: $?, waiting 5s..."; sleep 5' EXIT

export HOSTNAME=$HOSTNAME
export MAIL_HOSTNAME=ENTER HOSTNAME
export FROM_NUMBER=ENTER PLIVO FROM NUMBER TOKEN
export AUTH_ID=ENTER PLIVO AUTH TOKEN
export AUTH_TOKEN=ENTER PLIVO AUTH TOKEN
export MAIL_AUTH_PREFIX=gateway+SECRET

sed "s/__HOSTNAME__/$HOSTNAME/g; s/__MAIL_HOSTNAME__/$MAIL_HOSTNAME/g" > goguerrilla.conf <<EOF
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
  "GSMTP_LISTEN_INTERFACE": "0.0.0.0:25",
  "NGINX_AUTH": "127.0.0.1:8025",
  "NGINX_AUTH_ENABLED": "N",
  "SGID": "1000",
  "SUID": "1000"
}
EOF

echo "mails will be forwarded to +1234 if sent to gateway+SECRET-+1234@$MAIL_HOSTNAME"

rm -rf textmail
url=$(curl -s https://api.github.com/repos/cfstras/textmail/releases/latest \
	| grep browser_download_url | cut -d '"' -f 4)
echo "url: $url"
curl -SL -o textmail "$url"
#wget https://github.com/cfstras/textmail/releases/download/v1/textmail
chmod a+x textmail
exec ./textmail
