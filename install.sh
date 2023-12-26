#!/bin/bash

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
				echo ""
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
cd mdns
cp bin/mdns /usr/local/bin
cp mdns.service /etc/systemd/system/mdns.service
systemctl daemon-reload
systemctl enable mdns
mv /etc/systemd/resolved.conf /etc/systemd/resolved.bak
cp resolved.conf /etc/systemd/resolved.conf
ln -sf /run/systemd/resolve/resolv.conf /etc/resolv.conf
systemctl restart systemd-resolved
