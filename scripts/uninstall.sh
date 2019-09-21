#!/usr/bin/env bash
echo "Uninstalling goairmon server..."
systemctl stop goairmon
systemctl disable goairmon
rm -r /usr/local/goairmon || echo "Failed to delete app dir";
rm /etc/systemd/system/goairmon.service || echo "Failed to delete service";