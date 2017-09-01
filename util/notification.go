package util

import (
	"github.com/esiqveland/notify"
	"github.com/godbus/dbus"
)

func sendNotify(image, name, status string) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	n := notify.Notification{
		AppName:       "LightTower",
		ReplacesID:    uint32(0),
		AppIcon:       image,
		Summary:       name,
		Body:          status,
		ExpireTimeout: int32(5000),
	}
	if _, err := notify.SendNotification(conn, n); err != nil {
		return err
	}
	return nil
}
