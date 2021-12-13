#!/bin/bash

## Install golang app as systemd instanse
## Preconditions: root, pwd: [projname]
## Alertpoints: STDERR
## RUN FORM: approot NOT from system dir
## RUN: ./system/install_local-adverts.sh
apt-get update
apt-get -y install jpegoptim

APPNAME="local-adverts"
GO_USER="a"

# Deploy nginx configs
if [ -f /etc/nginx/nginx.conf-original ]; then
    cp -f ./system/nginx.conf /etc/nginx/nginx.conf
    echo "The special file 'nginx.conf' placed successfully in /etc/nginx"
else
    cp /etc/nginx/nginx.conf /etc/nginx/nginx.conf-original
    echo "The default file 'nginx.conf' backuped in /etc/nginx with -original suffix"
fi

# Backup default nginx and default-site configs
if [ -f /etc/nginx/sites-available/default-original ] || [ -f /etc/nginx/sites-available/default.conf-original ]; then
	echo "file default config backup file is available"
else
    unlink /etc/nginx/sites-enabled/default*;
    mv /etc/nginx/sites-available/default* /etc/nginx/sites-available/default-original;
    echo "The default 'default' file buckuped to -original suffix"
fi

if [ -f /etc/nginx/sites-available/$APPNAME ]; then
    echo "Nginx has previous $APPNAME.conf. ABORT INSTALLATION"
    exit 1
else
    cp -f ./system/$APPNAME".conf" /etc/nginx/sites-available/$APPNAME;
    ln -s /etc/nginx/sites-available/$APPNAME /etc/nginx/sites-enabled/$APPNAME
    # start nginx with new configs
    systemctl restart nginx
    echo "Sites config to nginx is installed"
fi

if [ -f /lib/systemd/system/$APPNAME.service ]; then
    echo "Systemd has previous $APPNAME.service"
    file /lib/systemd/system/$APPNAME.service
else
    cp -f ./system/$APPNAME".service" /lib/systemd/system
    chmod 644 /lib/systemd/system/$APPNAME".service"
    systemctl enable $APPNAME
    systemctl start $APPNAME
    systemctl status $APPNAME
    echo "Sytemd Unit of $SERVICEUNIT now installed: Ok!"
fi


echo "You will:"
echo "0. (make db) Import database data"
echo "1. edit configs/dev.yml prod file"
echo "2. run ./system/reloadcode.sh"
#systemctl daemon-reload
exit 0
