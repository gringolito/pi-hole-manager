#!/bin/sh

remove() {
    systemctl disable pi-hole-manager.service ||:
}

purge() {
    rm -rf /etc/pi-hole-manager
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
