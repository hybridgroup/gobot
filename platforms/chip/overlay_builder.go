package chip

import (
	"fmt"
	"os"
	"os/exec"
	"io/ioutil"
)

// BuildAndInstallOverlays uses the modified 'dtc' device tree compiler binary
// from https://github.com/atenart/dtc to build device tree overlay blobs and
// installs them for the overlay manager.
// Does not overwrite already existing dtbos.
// Requires root permissions.
func BuildAndInstallOverlays() (err error) {

	for _, info := range(overlays) {
		blobPath := overlayInstallPath + "/" + info.dtbo
		if _, err = os.Stat(blobPath); err == nil {
			// this blob is in place, check next
			continue
		}
		if os.IsNotExist(err) {
			tmpDir, err := ioutil.TempDir("", "dtbo")
			defer os.RemoveAll(tmpDir)

			err = buildOverlay(tmpDir, info.source)
			if err != nil {
				return fmt.Errorf("Failed to build overlay: %v", err)
			}
			err = installOverlay(tmpDir, info.dtbo)
			if err != nil {
				return fmt.Errorf("Failed to install overlay: %v", err)
			}
		} else {
			return fmt.Errorf("Failed to check for installed overlay: %v", err)
		}
	}

	return err
}

func buildOverlay(tmpDir string, source string) (err error) {
	path, err := exec.LookPath("dtc")
	if err != nil {
		return fmt.Errorf("Failed to find 'dtc' command on path")
	}

	sourcePath := tmpDir + "/overlay.dts"
	blobPath := tmpDir + "/overlay.dtbo"

	if err = ioutil.WriteFile(sourcePath, []byte(source), 0666); err != nil {
		return err
	}

	dtc := exec.Command(path, "-O", "dtb", "-o", blobPath, "-b", "o", "-@", sourcePath)
	if err = dtc.Run(); err != nil {
		return err
	}

	return nil
}

func installOverlay(tmpDir string, blobFile string) (err error) {
	if err = os.MkdirAll(overlayInstallPath, 0777); err != nil {
		return err
	}
	blobSource := tmpDir + "/overlay.dtbo"
	blobTarget := overlayInstallPath + "/" + blobFile
	err = copyFile(blobSource, blobTarget)
	return err
}
