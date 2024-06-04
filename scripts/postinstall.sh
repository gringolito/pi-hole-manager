#!/bin/sh
systemd_version=$(systemctl --version | head -1 | sed 's/systemd //g' | cut -d' ' -f 1)

install() {
    # RHEL/CentOS7 cannot use ExecStartPre=+ to specify the pre start should be run as root
    # even if you want your service to run as non root.
    if [ "${systemd_version}" -lt 231 ]; then
        printf "\033[31m  systemd version %s is less then 231, fixing the service file \033[0m\n" "${systemd_version}"
        sed -i "s/=+/=/g" /etc/systemd/system/dnsmasq-manager.service
    fi
    systemctl daemon-reload ||:
    systemctl unmask dnsmasq-manager.service ||:
    systemctl preset dnsmasq-manager.service ||:
    systemctl enable dnsmasq-manager.service ||:
    printf "\n"
    printf "  dnsmasq-manager service is enabled to start on the next system boot but it is not set to start automatically upon installation.\n"
    printf "\n"
    printf "  Please check the security configuration in:\n"
    printf "    /etc/dnsmasq-manager/config.yaml\n"
    printf "\n"
    printf "  and start the service with:\n"
    printf "    systemctl start dnsmasq-manager.service\n"
    printf "\n"
}

upgrade() {
    printf "\033[32m  Reloading dnsmasq-manager service\033[0m\n"
    systemctl try-restart dnsmasq-manager.service ||:
}

# Step 2, check if this is a clean install or an upgrade
action="$1"
if  [ "$1" = "configure" ] && [ -z "$2" ]; then
  # Alpine linux does not pass args, and deb passes $1=configure
  action="install"
elif [ "$1" = "configure" ] && [ -n "$2" ]; then
    # deb passes $1=configure $2=<current version>
    action="upgrade"
fi

case "$action" in
  "1" | "install")
    install
    ;;
  "2" | "upgrade")
    upgrade
    ;;
  *)
    # Alpine -> $1 == version being installed
    install
    ;;
esac
