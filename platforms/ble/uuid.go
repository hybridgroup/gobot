package ble

import (
	"fmt"
	"strconv"

	"tinygo.org/x/bluetooth"
)

func convertUUID(cUUID string) string {
	switch len(cUUID) {
	case 4:
		// convert to full uuid from "22bb"
		uid, e := strconv.ParseUint("0x"+cUUID, 0, 16)
		if e != nil {
			return ""
		}

		uuid := bluetooth.New16BitUUID(uint16(uid))
		return uuid.String()

	case 32:
		// convert "22bb746f2bbd75542d6f726568705327"
		// to "22bb746f-2bbd-7554-2d6f-726568705327"
		return fmt.Sprintf("%s-%s-%s-%s-%s", cUUID[:8], cUUID[8:12], cUUID[12:16], cUUID[16:20],
			cUUID[20:32])
	}

	return cUUID
}
