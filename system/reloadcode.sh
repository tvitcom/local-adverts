#!/bin/sh
# run as: ./system/reloadcode.sh
# remove old binary and copy new binary from dist then reload own service
APP_BIN="server"
PROJ_NAME="local-adverts"
DISTR_DIR="/build/";
PROJ_DIR="/home/user/Go/src/github.com/tvitcom/";
rm -f ${PROJ_DIR}${PROJ_NAME}/${APP_BIN};
cp -f ${PROJ_DIR}${PROJ_NAME}${DISTR_DIR}${APP_BIN} ${PROJ_DIR}${PROJ_NAME}/${APP_BIN};
systemctl restart ${PROJ_NAME};
systemctl status ${PROJ_NAME};
