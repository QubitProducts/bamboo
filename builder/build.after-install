#!/bin/bash
set -e

DATADIR="/var/bamboo"
OPT_DEFAULT_CONFIG_DIR="/opt/bamboo/config"
USER="bamboo"
GROUP="bamboo"
SHELL="/bin/false"

case "$1" in
	""|configure) # Work around https://github.com/jordansissel/fpm/issues/1272
        if [ ! -d "${DATADIR}/log" ]; then
			mkdir -p ${DATADIR}/log
		fi

        # Copy haproxy template without overrides
		cp -n ${OPT_DEFAULT_CONFIG_DIR}/haproxy_template.cfg ${DATADIR}
		# Copy example to production without overrides
		cp -n ${OPT_DEFAULT_CONFIG_DIR}/production.example.json ${DATADIR}/production.json

		if ! getent group | grep -q ${GROUP}; then
			groupadd -f ${GROUP}
		fi

		if ! getent passwd | grep -q ${USER}; then
			useradd -r -d ${DATADIR} --shell ${SHELL} -g ${GROUP} ${USER}
		fi
		chown -R ${USER}:${GROUP} ${DATADIR}

	;;

	abort-upgrade|abort-remove|abort-deconfigure)
	;;

	*)
		echo "postinst called with unknown argument \`$1'" >&2
		exit 1
	;;
esac
