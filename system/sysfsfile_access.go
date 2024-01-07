package system

import (
	"io"
	"os"
	"strconv"
	"strings"
)

type sysfsFileAccess struct {
	fs         filesystem
	readBufLen uint16
}

// Linux sysfs / GPIO specific sysfs docs.
//
//	https://www.kernel.org/doc/Documentation/filesystems/sysfs.txt
//	https://www.kernel.org/doc/Documentation/gpio/sysfs.txt
//	https://www.kernel.org/doc/Documentation/thermal/sysfs-api.txt
//	see also PWM.md
type sysfsFile struct {
	sfa       sysfsFileAccess
	file      File
	sysfsPath string
}

func (sfa sysfsFileAccess) readInteger(path string) (int, error) {
	sf, err := sfa.openRead(path)
	defer func() { _ = sf.close() }()
	if err != nil {
		return 0, err
	}

	return sf.readInteger()
}

func (sfa sysfsFileAccess) read(path string) ([]byte, error) {
	sf, err := sfa.openRead(path)
	defer func() { _ = sf.close() }()
	if err != nil {
		return nil, err
	}

	return sf.read()
}

func (sfa sysfsFileAccess) writeInteger(path string, val int) error {
	sf, err := sfa.openWrite(path)
	defer func() { _ = sf.close() }()
	if err != nil {
		return err
	}

	return sf.writeInteger(val)
}

func (sfa sysfsFileAccess) write(path string, data []byte) error {
	sf, err := sfa.openWrite(path)
	defer func() { _ = sf.close() }()
	if err != nil {
		return err
	}

	return sf.write(data)
}

func (sfa sysfsFileAccess) openRead(path string) (*sysfsFile, error) {
	f, err := sfa.fs.openFile(path, os.O_RDONLY, 0o644)
	if err != nil {
		return nil, err
	}
	return &sysfsFile{sfa: sfa, file: f, sysfsPath: path}, nil
}

func (sfa sysfsFileAccess) openWrite(path string) (*sysfsFile, error) {
	f, err := sfa.fs.openFile(path, os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	return &sysfsFile{sfa: sfa, file: f, sysfsPath: path}, nil
}

func (sfa sysfsFileAccess) openReadWrite(path string) (*sysfsFile, error) {
	f, err := sfa.fs.openFile(path, os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	return &sysfsFile{sfa: sfa, file: f, sysfsPath: path}, nil
}

func (sf *sysfsFile) close() error {
	if sf == nil || sf.file == nil {
		return nil
	}
	return sf.file.Close()
}

func (sf *sysfsFile) readInteger() (int, error) {
	buf, err := sf.read()
	if err != nil {
		return 0, err
	}

	if len(buf) == 0 {
		return 0, nil
	}

	return strconv.Atoi(strings.Split(string(buf), "\n")[0])
}

func (sf *sysfsFile) read() ([]byte, error) {
	// sysfs docs say:
	// > If userspace seeks back to zero or does a pread(2) with an offset of '0' the [..] method will
	// > be called again, rearmed, to fill the buffer.
	// > The buffer will always be PAGE_SIZE bytes in length. On i386, this is 4096.

	// TODO: Examine if seek is needed if full buffer is read from sysfs file.

	buf := make([]byte, sf.sfa.readBufLen)
	if _, err := sf.file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	i, err := sf.file.Read(buf)
	if i == 0 {
		return []byte{}, err
	}

	return buf[:i], err
}

func (sf *sysfsFile) writeInteger(val int) error {
	return sf.write([]byte(strconv.Itoa(val)))
}

func (sf *sysfsFile) write(data []byte) error {
	// sysfs docs say:
	// > When writing sysfs files, userspace processes should first read the
	// > entire file, modify the values it wishes to change, then write the
	// > entire buffer back.
	// however, this seems outdated/inaccurate (docs are from back in the Kernel BitKeeper days).

	// Write() returns already a non-nil error when n != len(b).
	_, err := sf.file.Write(data)
	return err
}
