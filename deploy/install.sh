#!/usr/bin/env bash

[ ! -n "$BASH_VERSION" ] && echo "You can only run this script with bash, not sh / dash." && exit 1

# Ebable bash strict mode
set -euxo pipefail

# Check docker version
if [ ! -x "$(command -v docker)" ]; then
    echo -e "Docker Engine not found.\n"
    echo "Installing Docker."
    dnf -y install dnf-plugins-core
    dnf config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo
    dnf install docker-ce docker-ce-cli containerd.io docker-compose-plugin
    systemctl start docker
fi

groupadd docker
usermod -aG docker borda

# Copy ssh keys
mkdir --parents --verbose /home/borda/.ssh
cp --verbose /root/.ssh/authorized_keys /home/borda/.ssh/

# Making base directory for borda
if [ ! -d /home/borda/.borda ]; then
    mkdir /home/borda/.borda
fi

chown --recursive --verbose borda:1000 /home/borda


echo "Starting Borda..."
curl --silent -SL https://github.com/gumrf/borda/blob/main/docker-compose.yaml -o /borda/docker-compose.yaml
cd /borda && docker compose up -d --force-recreate

echo -e "Congratulations! Borda is ready to use.\n"
echo "Visit http://$(curl -4s https://ifconfig.io):6900."