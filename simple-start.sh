#!/bin/bash
set -ex
trap 'echo "exit: $?, waiting 5s..."; sleep 5' EXIT

export MAIL_HOSTNAME=ENTER HOSTNAME
export FROM_NUMBER=ENTER PLIVO FROM NUMBER TOKEN
export AUTH_ID=ENTER PLIVO AUTH TOKEN
export AUTH_TOKEN=ENTER PLIVO AUTH TOKEN
export MAIL_AUTH_PREFIX=gateway+SECRET

# mails will be forwarded if sent to gateway+SECRET-+1123number@hostname

rm -rf textmail
url=$(curl -s https://api.github.com/repos/cfstras/textmail/releases/latest \
	| grep browser_download_url | cut -d '"' -f 4)
echo "url: $url"
curl -SL -o textmail "$url"
#wget https://github.com/cfstras/textmail/releases/download/v1/textmail
chmod a+x textmail
exec ./textmail
