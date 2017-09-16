package util

import (
	"github.com/esiqveland/notify"
	"github.com/godbus/dbus"
	"os"
)

func SendNotify(image, name, status, audio string) error {
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
		Hints:         map[string]dbus.Variant{},
	}
	if audio != "" {
		if _, err := os.Stat(audio); os.IsNotExist(err) {
			n.Hints["sound-name"] = dbus.MakeVariant(audio)
		} else {
			n.Hints["sound-file"] = dbus.MakeVariant(audio)
		}
	}
	if _, err := notify.SendNotification(conn, n); err != nil {
		return err
	}
	return nil
}
