package chip

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

const overlayInstallPath = "/lib/firmware/chip_io"
const overlayConfigPath = "/sys/kernel/config/device-tree/overlays"

type overlayInfo struct {
	dtbo   string
	folder string
	sysfs  string
}

var overlays = map[string]overlayInfo{
	"PWM0": {
		"chip-pwm0.dtbo", "chip-pwm", "/sys/class/pwm/pwmchip0",
	},
	"SPI2": {
		"chip-spi2.dtbo", "chip-spi", "/sys/class/spi_master/",
	},
}

func isLoaded(key string) (loaded bool, err error) {
	overlay, _ := keyToOverlay(key)
	configPath, _ := overlayToPaths(overlay)

	if _, err = os.Stat(configPath); err == nil {
		if _, err = os.Stat(overlay.sysfs); err == nil {
			return true, nil
		}
	}

	if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func overlayToPaths(overlay overlayInfo) (configPath string, overlayPath string) {
	configPath = overlayConfigPath + "/" + overlay.folder
	overlayPath = overlayInstallPath + "/" + overlay.dtbo
	return
}

func keyToOverlay(key string) (overlay overlayInfo, err error) {
	overlay, ok := overlays[key]
	if !ok {
		err = fmt.Errorf("Invalid overlay key %v", key)
	}
	return overlay, err
}

func copyFile(sourcePath string, destPath string) (err error) {
	cmd := exec.Command("cp", sourcePath, destPath)
	err = cmd.Run()
	return err
}

func LoadOverlay(key string) (err error) {
	overlay, err := keyToOverlay(key)
	if err != nil {
		return err
	}
	configPath, overlayPath := overlayToPaths(overlay)
	loaded, err := isLoaded(key)
	if err != nil {
		return err
	}
	if loaded {
		return fmt.Errorf("Overlay for %v already loaded!", key)
	}

	if err := os.MkdirAll(configPath, 0777); err != nil {
		return fmt.Errorf("Failed to create device tree path: %v", err)
	}

	err = copyFile(overlayPath, configPath+"/dtbo")
	if err != nil {
		return err
	}

	time.Sleep(200 * time.Millisecond)

	loaded, err = isLoaded(key)
	if err != nil {
		return err
	}
	if !loaded {
		return fmt.Errorf("Failed to load overlay for %v", key)
	}

	return nil
}

func UnloadOverlay(key string) (err error) {
	overlay, err := keyToOverlay(key)
	if err != nil {
		return err
	}
	configPath, _ := overlayToPaths(overlay)
	err = os.Remove(configPath)
	return err
}
