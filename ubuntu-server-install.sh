#!/bin/bash

echo "Installing mDNS"

	if [ "$EUID" -ne 0 ]; then
    echo "Sorry, you need to run this as root"
		exit 1
	fi

if [[ -e /etc/debian_version ]]; then
		OS="debian"
		source /etc/os-release

if [[ $ID == "ubuntu" ]]; then
			OS="ubuntu"
			MAJOR_UBUNTU_VERSION=$(echo "$VERSION_ID" | cut -d '.' -f1)
			if [[ $MAJOR_UBUNTU_VERSION -lt 22 ]]; then
				echo "If you're using Ubuntu < 22.04, then you can continue, at your own risk."
				until [[ $CONTINUE =~ (y|n) ]]; do
					read -rp "Continue? [y/n]: " -e CONTINUE
				done
				if [[ $CONTINUE == "n" ]]; then
					exit 1
				fi
			fi
        fi
else
        echo "Looks like you aren't running this installer on Ubuntu system"
		exit 1
fi

cd /tmp
git clone https://github.com/MarlikAlmighty/mdns
ERR=$?
if [[ $ERR != 0 ]]; then
echo "error while cloning mDNS from github, exit."
	exit 1
fi

cd mdns
cp bin/mdns /usr/local/bin
ERR=$?
if [[ $ERR != 0 ]]; then
	echo "error while installing mdns, exit."
	exit 1
fi

chmod 755 /usr/local/bin/mdns
ERR=$?
if [[ $ERR != 0 ]]; then
	echo "error while chmod /usr/local/bin/mdns, exit."
	exit 1
fi

cp mdns.service /etc/systemd/system/mdns.service
ERR=$?
if [[ $ERR != 0 ]]; then
	echo "error while installing mdns.service, exit."
	exit 1
fi

systemctl daemon-reload
ERR=$?
if [[ $ERR != 0 ]]; then
echo "error while daemon-reload, exit."
exit 1
fi

systemctl enable mdns
ERR=$?
if [[ $ERR != 0 ]]; then
echo "error while enable mdns, exit."
exit 1
fi

mv /etc/systemd/resolved.conf /etc/systemd/resolved.bak
ERR=$?
if [[ $ERR != 0 ]]; then
echo "error while mv resolved.conf, exit."
exit 1
fi

cp resolved.conf /etc/systemd/resolved.conf
ERR=$?
if [[ $ERR != 0 ]]; then
echo "error while cp resolved.conf, exit."
exit 1
fi

ln -sf /run/systemd/resolve/resolv.conf /etc/resolv.conf
ERR=$?
if [[ $ERR != 0 ]]; then
echo "error while ln resolv.conf, exit."
exit 1
fi

systemctl restart systemd-resolved
ERR=$?
if [[ $ERR != 0 ]]; then
echo "error while restart systemd-resolved, exit."
exit 1
fi

ufw allow 53/tcp
ufw allow 53/udp
ERR=$?
if [[ $ERR != 0 ]]; then
echo "here is error: ufw allow 53 ports, do it manually"
fi

systemctl start mdns
ERR=$?
if [[ $ERR != 0 ]]; then
echo "error while start mdns, exit."
exit 1
fi

rm -rf /tmp/mdns
ERR=$?
if [[ $ERR != 0 ]]; then
echo "error while rm -rf /tmp/mdns, exit."
exit 1
fi

echo "Done, mdns is installed."

exit 0