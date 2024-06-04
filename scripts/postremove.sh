#!/bin/sh

remove() {
    systemctl disable dnsmasq-manager.service ||:
}

purge() {
    rm -rf /etc/dnsmasq-manager
}

upgrade() {
    true
}

action="$1"

case "$action" in
  "0" | "remove")
    remove
    ;;
  "1" | "upgrade")
    upgrade
    ;;
  "purge")
    purge
    ;;
  *)
    # Alpine
    remove
    ;;
esac
