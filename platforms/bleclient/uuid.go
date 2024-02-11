package bleclient

import (
	"fmt"
	"strconv"
	"strings"

	"tinygo.org/x/bluetooth"
)

// convertUUID creates a common 128 bit UUID xxxxyyyy-0000-1000-8000-00805f9b34fb from a short 16 bit UUID by replacing
// the yyyy fields. If the given ID is still an arbitrary long one but without dashes, the dashes will be added.
// Additionally some simple checks for the resulting UUID will be done.
func convertUUID(cUUID string) (string, error) {
	var uuid string
	switch len(cUUID) {
	case 4:
		uid, err := strconv.ParseUint(cUUID, 16, 16)
		if err != nil {
			return "", fmt.Errorf("'%s' is not a valid 16-bit Bluetooth UUID: %v", cUUID, err)
		}
		return bluetooth.New16BitUUID(uint16(uid)).String(), nil
	case 32:
		// convert "22bb746f2bbd75542d6f726568705327" to "22bb746f-2bbd-7554-2d6f-726568705327"
		uuid = fmt.Sprintf("%s-%s-%s-%s-%s", cUUID[:8], cUUID[8:12], cUUID[12:16], cUUID[16:20], cUUID[20:])
	case 36:
		uuid = cUUID
	}

	if uuid != "" {
		id := strings.ReplaceAll(uuid, "-", "")
		_, errHigh := strconv.ParseUint(id[:16], 16, 64)
		_, errLow := strconv.ParseUint(id[16:], 16, 64)
		if errHigh == nil && errLow == nil {
			return uuid, nil
		}
	}

	return "", fmt.Errorf("'%s' is not a valid 128-bit Bluetooth UUID", cUUID)
}
