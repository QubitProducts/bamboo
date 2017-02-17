#!/bin/bash

DATADIR="/var/bamboo"
USER="bamboo"
GROUP="bamboo"
SHELL="/bin/false"

set -e

if getent passwd | grep -q "^${USER}:"; then
	userdel ${USER}
fi

if getent group | grep -q "^${GROUP}:"; then
	groupdel ${GROUP}
fi
