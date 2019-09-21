#!/bin/bash

if [[ ! -f goairmon ]] || [[ ! -d resources ]]; then
	echo "This script should be run from the unzipped dist file directory"
	exit 1
fi

read -p "Set Username: " UserName
read -s -p "Password: " Password
echo

AppDir=/usr/local/goairmon

mkdir ${AppDir} || echo "Failed to make root dir"
cp -r ./. ${AppDir}
cp goairmon.service /etc/systemd/system/goairmon.service

rm ${AppDir}/install.sh ${AppDir}/goairmon.service

(cd ${AppDir} && cmd/adduser -username ${UserName} -password ${Password}) || echo "Failed to add new user"

CookieKey=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 16 | head -n 1)
sed -i "s/%%COOKIE_KEY%%/${CookieKey}/g" ${AppDir}/.env

setcap 'cap_net_bind_service=+ep' ${AppDir}/goairmon;

systemctl enable goairmon;
systemctl start goairmon;

systemctl is-active --quiet goairmon;
if [[ $? -eq 0 ]]; then
		echo "Goairmon installed successfully";
else
		echo "Failed to start goairmon service";
fi
