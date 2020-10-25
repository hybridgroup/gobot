package ble

import (
	"fmt"
)

func convertUUID(cUUID string) string {
	switch len(cUUID) {
	case 4:
		// 2a270000-0000-0000-0000-000000000000
		// convert "22bb"
		// to "22bb0000-0000-0000-0000-000000000000"
		return fmt.Sprintf("%s0000-0000-0000-0000-000000000000", cUUID)
	case 32:
		// convert "22bb746f2bbd75542d6f726568705327"
		// to "22bb746f-2bbd-7554-2d6f-726568705327"
		return fmt.Sprintf("%s-%s-%s-%s-%s", cUUID[:8], cUUID[8:12], cUUID[12:16], cUUID[16:20],
			cUUID[20:32])
	}

	return cUUID
}
