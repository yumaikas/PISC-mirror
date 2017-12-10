// Code generated by go-bindata.
// sources:
// stdlib/bools.pisc
// stdlib/debug.pisc
// stdlib/dicts.pisc
// stdlib/io.pisc
// stdlib/locals.pisc
// stdlib/loops.pisc
// stdlib/math.pisc
// stdlib/random.pisc
// stdlib/shell.pisc
// stdlib/std_lib.pisc
// stdlib/strings.pisc
// stdlib/symbols.pisc
// stdlib/vectors.pisc
// stdlib/with.pisc
// DO NOT EDIT!

package pisc

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _stdlibBoolsPisc = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\xd7\xe2\x4a\xca\xcf\xcf\x29\xd6\x2b\xc8\x2c\x4e\xe6\xd2\xd2\xe7\xe2\x52\x56\xf0\x54\x48\x2f\x4d\x2d\x2e\x56\x28\xc9\xc8\x2c\x56\xc8\x2c\x56\x48\xcd\x2d\x28\xa9\x54\xc8\xcb\x2f\xb7\xe7\xe2\x02\x04\x00\x00\xff\xff\xcb\x3e\x17\xc6\x30\x00\x00\x00")

func stdlibBoolsPiscBytes() ([]byte, error) {
	return bindataRead(
		_stdlibBoolsPisc,
		"stdlib/bools.pisc",
	)
}

func stdlibBoolsPisc() (*asset, error) {
	bytes, err := stdlibBoolsPiscBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "stdlib/bools.pisc", size: 48, mode: os.FileMode(436), modTime: time.Unix(1511943077, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _stdlibDebugPisc = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\xcc\xb1\x0a\xc2\x30\x10\xc6\xf1\xbd\x4f\xf1\xd1\x49\x87\xbc\x40\x8b\x2e\xba\x3b\x38\x8a\x43\xb8\x9e\x26\x10\x72\xe1\x2e\x1a\x7d\x7b\x91\x2a\x74\xe8\xf6\xc1\x9f\xef\xd7\x0d\xc7\xd3\x01\x16\xa4\xb9\xa2\x7c\x8b\x2f\xd7\x44\x27\xc3\x06\xce\x61\x8b\x73\x90\x66\x68\x21\x52\xc0\xdc\x31\x77\xaf\x0c\x7a\xa8\x72\xae\xe9\x8d\x24\x7e\xe2\x09\x63\xb7\xe0\x92\x90\x4f\xab\xd0\xaf\x7c\x89\x98\xff\x0a\x8c\xa4\x30\xc6\x6e\x58\xbb\x5f\x00\x79\xb2\x62\x6f\x55\x63\xbe\xa3\xc7\x0e\x3d\xac\xf9\x02\xab\xea\x48\x32\xf9\xba\x9c\x39\x16\x14\x8d\xb9\xe2\x0a\xf6\x14\x66\x0f\x23\x3e\x01\x00\x00\xff\xff\xd4\x33\x23\xc6\xf1\x00\x00\x00")

func stdlibDebugPiscBytes() ([]byte, error) {
	return bindataRead(
		_stdlibDebugPisc,
		"stdlib/debug.pisc",
	)
}

func stdlibDebugPisc() (*asset, error) {
	bytes, err := stdlibDebugPiscBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "stdlib/debug.pisc", size: 241, mode: os.FileMode(436), modTime: time.Unix(1500942096, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _stdlibDictsPisc = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x94\xcd\x6e\xdb\x3a\x10\x85\xf7\x7a\x8a\x03\x43\x17\xb8\x16\x4a\x69\x2f\x1b\xf1\xa6\xed\xba\x68\x97\x86\x17\x8c\x34\xb2\x88\xc8\xa4\x4a\x8e\x1c\xf8\xed\x0b\x52\x3f\x96\x65\xa5\x69\x36\x49\x86\x67\x66\xbe\x39\x1c\x31\x4b\xf0\x55\x15\xac\x8c\x96\x56\x91\x83\xd4\x25\xb8\x26\x65\xd1\xb1\x6a\x14\xfb\x58\x92\x45\x51\x8e\x52\x15\x2c\x54\x25\xe8\xd2\xf2\x4d\x38\x96\xc5\x1b\xfe\x47\x9a\x42\x88\x70\x76\xc0\x16\x21\xda\x2b\x0e\x38\x62\xef\xe3\x2f\x38\xe1\xbd\x26\x8d\xdd\xac\x88\x36\x2c\xfc\xdf\x0f\x15\xb0\x45\xd9\xb5\xe0\x5b\x4b\xa6\xc2\x66\xe2\xba\x6d\xe0\xd8\x0a\x4d\xbf\xd7\x8b\x92\x76\x9d\xa5\x50\xaf\xd7\x3f\x55\x5d\x41\x7f\x22\xf1\x95\xce\xc4\xc2\x58\x51\x52\x25\xbb\xc6\xc3\x85\x93\x37\xba\x61\x0c\x09\x81\xab\x6c\x3a\xc2\x16\xf9\x18\xcb\xbd\x20\xf7\xd2\x08\xc3\x4f\x1c\x12\xe3\x90\xe9\x1b\xd5\xd2\x89\x37\x0a\xa6\x2c\x8f\xce\xc4\x38\xf9\xf0\x50\xed\x04\x55\x45\xbb\x28\xca\x12\x7c\x27\x2e\x6a\x54\xd6\x5c\x30\x1b\x2e\xc9\x22\xe4\x3f\x7e\x7e\x83\x78\x99\x13\xf6\x64\xe3\xb4\xbe\xea\x6e\xd2\xad\x09\x07\x6f\x8e\xc1\xf3\xbc\xc4\x09\xa5\x6a\xef\xc9\x71\x89\x1e\xe2\x17\xb1\xdf\x88\x90\x5a\x19\x0b\x89\xb3\xba\x92\xee\x89\xb4\xf1\x44\xc1\x91\x2f\x78\xed\x18\x0d\xc9\x2b\x05\xfd\x8c\xd8\xe8\x10\xe9\x8d\x9f\xf0\xf7\x7b\x31\x62\x79\x9e\x01\x6d\x7e\x65\x6d\xe7\xea\x69\x8a\x55\xf5\x20\x74\x61\xda\x69\xdc\x78\x31\xee\x36\x5c\x6c\x63\x8a\x61\xec\xa5\x3b\x4b\xfd\x83\x35\xbd\x2f\x2b\x05\xdc\xbb\x6c\xef\x5d\xf7\x22\xbe\xf3\x59\x2d\x2f\xb4\xda\xd9\xcd\x3a\xef\xff\x31\x67\xb0\x61\xc8\x8a\x3f\x30\x22\xec\xde\xc6\x9f\xd6\xd2\xe1\x95\xfc\x15\x51\x6b\xa9\x90\x4c\x65\x8a\xce\x11\xe2\xab\xb4\x1e\xd4\xff\x52\xda\x31\xc9\x72\x03\xb2\xd6\x58\x64\x49\x18\xf6\xde\x39\x4c\x37\x21\x27\xd9\x44\x1d\xff\xf5\xde\x46\x8a\xcf\x31\x3e\xe4\x38\x3e\x18\x3f\xac\xe5\x67\x6c\x7e\x53\xc3\xaa\x4a\x0d\xd2\x6c\x6f\x60\x03\xb6\x1d\x85\x9d\x25\x59\xd4\xe8\xb4\x2a\x4c\x49\x28\x6a\x69\x65\xc1\xe4\x7b\x83\x6b\xe5\xfc\xeb\xa2\xf4\x79\xf6\x61\xfd\x37\x0c\x78\x78\xda\x88\xe7\xb7\xa6\x07\x3c\xf6\x0f\xc0\xf4\x5f\xff\x95\x07\x50\x5e\xf0\x9e\x02\x8e\xf0\x18\x83\x6c\xb6\xba\x87\xb5\x55\x3c\x8c\x7b\x3e\x3d\x23\x53\x46\xba\xd0\xa7\xe9\xfc\x05\x28\x64\xd3\x60\x17\xfd\x09\x00\x00\xff\xff\x7e\x91\x17\xc6\xe5\x05\x00\x00")

func stdlibDictsPiscBytes() ([]byte, error) {
	return bindataRead(
		_stdlibDictsPisc,
		"stdlib/dicts.pisc",
	)
}

func stdlibDictsPisc() (*asset, error) {
	bytes, err := stdlibDictsPiscBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "stdlib/dicts.pisc", size: 1509, mode: os.FileMode(436), modTime: time.Unix(1512944015, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _stdlibIoPisc = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x54\x41\x6f\xf3\x36\x0c\x3d\xd7\xbf\xe2\x01\x19\x90\xe6\x43\x9c\xef\x3a\xb4\xc0\x86\x21\x6d\x80\x5e\x96\x61\xeb\xb0\x4b\x81\x55\xb6\xe9\x58\x83\x22\x79\x12\xd5\xc0\xff\x7e\x20\x6d\x27\x59\xbb\x0e\xdf\x4d\x96\xc8\xf7\x1e\xc9\x47\x7f\xfd\x82\xe7\x8e\xf0\xb4\x87\xb3\x55\x34\xd1\x52\x42\x68\xf1\xcb\xd3\x6f\x5b\x7c\xf9\x5a\x14\xfa\x6e\x13\x9a\x50\xc3\x26\xbf\x64\xb4\x21\xc2\x20\xf5\x54\xdb\xd6\xd6\x38\x85\xd8\x48\xe4\xdd\xc3\x7e\x2b\x30\xb7\x28\x4b\xac\x50\x2c\xf0\x2b\x99\x86\x62\x2a\x8a\x6d\x8e\x91\x3c\xbb\x61\x3d\xe2\x76\x26\x21\x92\x69\xac\x3f\xc0\xf8\x06\xa7\x68\x59\xce\xa1\xfa\x8b\x6a\x4e\xe0\xce\x30\x4c\x24\x54\x21\xfb\x06\x1c\x90\x7d\x43\xd1\x0d\x12\x74\x08\xa8\x4c\xa2\x06\x55\x6e\x6d\x98\x73\x36\x78\x34\x75\xa7\xa8\x14\x95\x80\x3b\x42\x1b\x9c\x0b\x27\xc9\x6a\xb3\xaf\xd9\x06\x9f\xe0\x82\x69\xa8\x81\xf5\x1c\x60\x39\xa1\xb1\xfa\x60\xe2\x00\xf2\x2c\xf5\xdf\x15\xc5\x62\xa1\x50\x65\x35\x30\x8d\x15\xe9\x69\x55\x48\x4d\x09\x66\xfc\x6c\x63\x38\x2a\xcf\x44\x6b\x12\x8c\x17\xe4\x35\x12\xb1\x96\x74\x79\x4d\x78\xdc\xef\xa4\x16\x8e\x99\x60\x5b\x7d\x22\xdf\x48\xb7\x55\xab\x75\x04\xab\x8d\xa9\x3b\x6a\x2e\x1a\x62\xf6\x93\x06\x3d\x5d\x34\xfc\xfe\xbc\x2b\xbf\x1f\x2f\xdf\x29\xb9\x24\x3b\x3b\x27\x27\x8e\x57\xb9\x7a\xff\x2e\x6b\x0d\x8e\xf6\x78\x14\xd9\xcb\x97\xb8\xd4\xd1\x2c\x5f\xfc\x12\xa1\x3d\xab\x55\x64\x29\x44\x31\x1f\xf7\xbb\x1f\x15\x94\x73\xf4\xe9\x7f\x2b\x93\x91\x54\x44\x7e\xae\x6f\x8d\xd6\xb8\x44\x08\xdc\x51\x3c\xd9\x44\x45\x71\x5f\x8c\x26\x0a\x3d\xf9\x52\x92\xca\xa9\xaf\xb7\xe8\x0d\x77\x42\x38\x5a\x0a\x2b\xec\x04\x73\x6e\x6c\x67\xde\x08\x06\xaf\xb5\x0b\x89\x5e\x75\x8c\xc3\x68\xa2\xda\x78\x54\x84\xda\x38\x47\x6a\x24\x0d\x51\x55\x57\x96\x9a\x04\xfa\xc6\x11\x3e\x8a\x10\x77\xfe\x4b\xc4\x74\xb1\xc2\x1f\x7a\x48\x6a\xd5\x64\x8f\xbd\xa3\x28\xb4\x7e\x16\xb6\xc6\xe8\xc3\x41\x75\x04\xef\x86\x31\x57\xbb\x9f\x10\xa2\x0c\xc5\xfa\x43\xc2\x7d\x71\x33\xd2\xea\x7b\x39\x5e\xe3\x56\x87\xa6\xeb\xa4\x54\xfa\x29\x0e\xea\x08\x86\x59\xfb\x38\xab\x99\x75\x8f\x00\xd3\xd8\xbf\x39\x7d\x0d\xd3\xf7\xe4\x75\x23\x5f\xfc\x1c\x24\x33\x9c\x71\xa5\x15\xd2\x80\x1f\xce\xda\xe4\xe6\x67\x73\xa4\xc9\x5c\xdb\xe0\x99\x3c\x27\xac\x80\x9f\x90\xd9\x3a\xcb\xc3\x79\xef\xf4\xaf\x31\xef\x7c\x3a\x1a\xe7\x34\x3f\x8d\x7b\x38\xb5\x61\x83\x87\x00\x1f\x18\x59\x8c\xe1\xe1\x4c\x3c\xd0\x14\x77\x16\xd2\x47\xfb\xf6\x67\x9f\x39\x5d\xd7\x77\x75\x6b\xd3\xfb\xf1\xea\x1f\x2a\xcb\x2f\x43\x44\xf4\xd1\x7a\x5d\xcd\x37\xe3\xf2\x0c\x8c\x93\xe5\xae\x0c\x99\xfb\xcc\xf3\xa4\xff\xce\x81\x05\x7d\xb3\xc1\xaa\xb8\xb9\x93\xcf\xe2\xe6\x83\x2d\x9a\xdc\xe3\x3b\x0d\x15\x08\x6c\xd4\x60\x62\xe5\x05\x9e\xf7\x0f\xfb\x62\x81\x09\x7c\x6c\xb0\x70\xe0\x4c\x68\xfd\x27\x7c\xf8\x0f\xc2\x69\x19\x3e\x21\xfc\x27\x00\x00\xff\xff\xb2\x42\xc8\x50\xc8\x05\x00\x00")

func stdlibIoPiscBytes() ([]byte, error) {
	return bindataRead(
		_stdlibIoPisc,
		"stdlib/io.pisc",
	)
}

func stdlibIoPisc() (*asset, error) {
	bytes, err := stdlibIoPiscBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "stdlib/io.pisc", size: 1480, mode: os.FileMode(436), modTime: time.Unix(1491353691, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _stdlibLocalsPisc = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x92\x4f\x8b\xdb\x30\x10\xc5\xef\xfe\x14\x0f\xb2\x87\x75\x16\x27\x77\x07\x0a\x25\x0d\x65\xa1\xb4\x61\x73\x0c\xa1\x0c\xf2\x38\x11\x49\x24\x57\x1a\x67\xfb\xf1\x8b\xfe\x38\xeb\xac\x7b\xd0\x41\xcc\x9b\xdf\xbc\xf9\xb3\x9c\xe3\x87\x55\x74\xf1\xd0\x06\xdb\xd7\xdd\x1a\xf3\x65\x51\xd4\xdf\x7e\xad\xe1\x59\xaa\x4b\x88\xe1\x19\x37\xba\xc0\xd0\x95\x51\x55\x28\xb1\x63\x49\x3f\xb1\x31\x82\x55\x4e\x39\x8e\x52\x06\x79\x10\x94\xf8\xce\x02\x39\x71\xf8\xf5\x0c\xf2\xde\x2a\x4d\xc2\x0d\xde\xb5\x9c\x40\x48\x59\x29\x3a\xc5\x79\x3c\xa7\xca\xdb\xde\x07\xb5\x17\x52\x67\xb4\x2e\x94\x68\xad\x43\x16\x59\x23\x36\x56\xc9\xff\x24\x1b\x68\x8d\xb3\xdd\x67\x9c\xed\x40\x99\x63\xdb\x3b\xa6\x6d\xa7\x94\x3b\x86\x49\x9d\xee\x4d\xfe\xe9\xad\x60\x8f\x33\x6e\xa8\xb0\x58\xe0\x10\xb8\x29\x6f\xb1\x40\x09\xbc\xf5\x06\x14\x65\x24\xda\x9a\xe8\x36\x10\x72\xc3\xda\xc4\x4a\xaa\x77\x8e\x8d\x3c\xf4\xb5\xdb\x6c\xea\x98\xf8\xa5\xd1\x4a\x62\xfd\xe5\x1c\x6b\xdb\x69\x6e\xd0\x3a\x7b\x45\x4b\x4a\xac\x0b\x0b\xab\xa1\x4e\x64\x8e\x3c\x38\xba\x91\xfb\x99\xc7\x1f\x6d\xf8\x77\xea\xb0\xc7\x7e\xb4\xa0\x03\xce\xcc\x1d\x0e\x68\x74\x17\xdf\xc7\xba\x57\x28\x8a\x19\xbe\x8a\xf0\xb5\x13\x6d\x8e\x61\xcd\x99\x9f\x17\x55\xd4\xdb\xb7\x0d\x9e\xea\xff\xd4\x2b\x07\x69\xf2\x1b\x8f\x0b\x5b\xc7\xad\xfe\xcb\x3e\x78\x9d\xe1\xd5\x28\xc7\xd7\xd0\x2f\x99\x06\x0d\xe7\x5f\x81\x48\xad\xaa\xd1\xf1\x94\x31\x9c\x7c\x55\x37\x72\x58\x65\xd5\xcb\xcb\x83\x4a\x9b\x4f\xaa\x62\x16\x4f\x2e\x14\xf0\x2c\xd9\x75\x4a\x7d\x9a\x1c\xe7\xc7\x54\x06\x7c\x3d\xb9\xf9\xd1\x78\x8a\x2c\xfa\x1d\x54\x71\x39\x8f\xb8\xa6\xef\x46\xc8\x38\xfb\xa0\xaa\x82\x91\xd5\xbf\x00\x00\x00\xff\xff\xaa\x09\x9f\x67\x71\x03\x00\x00")

func stdlibLocalsPiscBytes() ([]byte, error) {
	return bindataRead(
		_stdlibLocalsPisc,
		"stdlib/locals.pisc",
	)
}

func stdlibLocalsPisc() (*asset, error) {
	bytes, err := stdlibLocalsPiscBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "stdlib/locals.pisc", size: 881, mode: os.FileMode(436), modTime: time.Unix(1504845781, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _stdlibLoopsPisc = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\xcb\x31\x0a\x02\x31\x10\x46\xe1\x7e\x4e\xf1\x97\xeb\xc2\x66\x7b\x53\xaa\x9d\xe0\x15\x46\xe2\x80\x81\x98\x44\x67\xa2\x1e\x5f\x42\xac\xb7\x7d\x7c\x6f\x9d\xe9\x5c\x4a\x55\x57\xa3\x06\x9a\x57\xa2\xfd\xf1\x72\x80\xc5\x87\x28\x26\x64\x3c\x5b\x31\x2c\x0b\x9c\xc3\x0e\xa7\xaf\x84\x66\x02\xee\x95\x91\xff\xce\x8f\xe9\x73\x8f\x49\x30\xa1\xbe\xe4\xb6\xf9\x0d\xc8\xdd\x31\x92\x5c\xdf\xa2\x60\x63\xf8\x5f\x00\x00\x00\xff\xff\x8a\xf6\x37\x23\x8e\x00\x00\x00")

func stdlibLoopsPiscBytes() ([]byte, error) {
	return bindataRead(
		_stdlibLoopsPisc,
		"stdlib/loops.pisc",
	)
}

func stdlibLoopsPisc() (*asset, error) {
	bytes, err := stdlibLoopsPiscBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "stdlib/loops.pisc", size: 142, mode: os.FileMode(436), modTime: time.Unix(1491353691, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _stdlibMathPisc = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x95\x4f\x6f\xe3\x36\x10\xc5\xef\xfa\x14\xef\xb8\x71\x23\x24\xca\xd1\x71\x1c\x74\xbd\xdd\xb6\xc0\xa2\x3e\x24\xb7\x20\x80\x29\x69\x24\x33\x95\x48\xed\x90\x8a\xed\x7e\xfa\x82\xb4\x6c\x51\xf2\x1f\xec\xcd\x30\x7f\xef\x0d\xf9\x38\x43\x4d\xbf\x2d\x17\x98\xe7\xba\x4d\x2b\xc2\x17\x6c\x11\xc7\xc8\x71\x83\x85\x26\xce\x08\x02\xaa\xad\x53\x62\x58\x0d\x81\x0e\x7b\x8c\xbc\x4a\x64\xda\x1c\x24\x3b\xdc\xe0\x77\xce\x90\x69\x23\x15\xe1\x11\x3d\xb2\x0e\x99\xbf\x76\x0d\x71\xaa\x2b\x99\x41\x9c\xc1\x8d\x54\x63\xc7\x13\xe0\x9a\xdf\x10\xb6\xe2\xc4\xcd\x0a\x55\x92\xb2\x03\xe6\x9a\xa1\xe7\x7b\x3c\x4b\xd9\x86\xf4\xa2\x4d\x09\xac\x75\x40\x90\xac\x06\x04\xc9\x4a\xaa\x12\xba\xe8\xf3\x6b\x98\x32\x69\xa4\x56\x87\x70\x8f\xe2\x61\xa0\x8b\x61\x3a\x57\xb2\x1c\xe5\x48\x5c\x84\xe0\x1f\xcc\x9a\x51\xb4\x2a\xb3\xae\x68\x40\x65\xc3\x72\x75\x53\x51\x4d\xca\x0a\xde\x81\x2e\x88\xb6\xcd\xc0\x7a\xdb\x68\x45\xca\x4a\x31\xa6\x1e\x42\xec\x75\xa3\x61\xd7\x84\x46\x6f\x88\x5d\x16\xab\xed\x2a\x84\xeb\x24\xa4\x57\xae\x48\x82\x78\x75\x8b\x5a\xf3\x21\x2f\x42\xa1\x19\xa6\x16\x55\xd5\xe5\x66\x70\xf4\x28\x2a\xad\x39\xf4\xf8\xee\xff\xf8\xa5\xd4\x4b\x51\xd7\x22\xd4\xfe\xe9\xff\x38\x1e\xfd\x75\xf9\x6d\x79\x84\x3f\xee\x43\xf2\xe3\x1e\x5f\xc9\x18\xaa\x4e\x83\xfa\x18\x1c\xe9\xbb\x64\x63\xf1\xaf\x54\xf9\x45\x41\xa5\xcb\x50\xf1\x8f\xb0\x2d\x8b\x0a\x3f\x74\x29\x58\xda\x75\x8d\x47\xf4\x64\x32\xd8\xc6\x0f\x5d\x22\x15\x86\x90\xdc\x87\x76\xc9\xe0\xaa\x56\x5b\x24\xf8\xcd\x97\x59\xdd\x22\x6d\xed\x3e\x5c\x91\x65\x2d\x0b\x4b\xd8\x48\xbb\xee\xe2\xfd\x14\x55\x4b\x26\xb4\x7a\x38\x5b\xee\x21\x44\xd2\x10\xf9\x2a\x95\x6f\xa2\xae\x3d\xdc\x45\x6c\x8f\xf0\x68\xca\x5f\x5c\xeb\x9e\xc4\x71\x65\xd2\xcd\x79\xc1\xcf\xe1\x6c\xbe\xfc\x6c\x05\x77\xd3\x79\x02\x8f\x9e\x86\xd7\xee\x59\x38\xc7\x5d\xda\x85\xbd\xa8\xe1\x56\x0d\x06\xeb\x6f\x65\xa9\x24\xde\xc7\xea\xa3\xb8\x85\x30\xc1\x7b\xda\x09\x77\x83\x4b\x5d\x72\x4e\x1c\xff\x47\xac\x4f\x5a\x46\x17\x7e\x9a\x0c\x65\x5a\xe5\xfb\xae\x3a\x7a\x24\xa7\x1e\x5a\xd1\x2f\x5a\xdc\x4d\x22\xd7\xed\x53\xb8\x2e\x70\xeb\x76\xa3\x21\xb8\x44\x2d\xec\xfa\x28\x36\xd1\xe4\x2e\x9a\xa2\xf6\xd7\x28\x90\xba\x6a\xbe\x71\x88\x71\x83\x87\xbc\x6d\x30\xc7\x1b\x94\x6c\xf0\x8e\x37\xe4\xac\xdd\x0f\x59\xb8\x0f\x07\x6a\xb1\xed\x55\x95\xe0\xd2\x8b\xd0\xcb\x3a\xfc\xa0\xf7\xb2\xe8\x6e\x82\x17\x5d\x93\x6b\x3b\x99\xed\x77\xe3\xf7\x40\x9f\xa4\x9e\xf1\x05\xca\xb9\x3d\xbb\xea\xa8\x75\x0e\x97\xda\xb3\x2f\x97\xcb\x4f\x69\x34\xef\x99\xfa\x40\x85\x4c\x34\xc5\x53\xbf\x23\xb7\x1a\xf7\x6b\x77\x13\xcc\x20\x0d\x44\xc5\x24\xf2\x1d\x72\x2a\xa4\xa2\x7c\x5f\x7c\x3e\xd2\xcd\xa0\xdc\x97\x20\x9a\x62\x36\x5a\xf1\x87\x9b\xe1\x0d\x4f\x78\x47\x2e\x1b\x68\xf6\xdc\x7c\x64\xf0\xd4\x39\x44\x53\x88\xd4\xf8\xc5\x38\x86\xc0\x0d\x9c\xc1\xbd\xb7\x88\x13\x4c\xf0\x8e\xcd\x9a\x5c\xcb\xfd\x1f\x00\x00\xff\xff\x5a\xd8\x7d\x5b\xc0\x07\x00\x00")

func stdlibMathPiscBytes() ([]byte, error) {
	return bindataRead(
		_stdlibMathPisc,
		"stdlib/math.pisc",
	)
}

func stdlibMathPisc() (*asset, error) {
	bytes, err := stdlibMathPiscBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "stdlib/math.pisc", size: 1984, mode: os.FileMode(436), modTime: time.Unix(1500942096, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _stdlibRandomPisc = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x14\xca\xc1\x09\x80\x30\x0c\x40\xd1\x55\xfe\x51\x0f\x59\x40\xa7\x09\x49\xc0\x42\x9b\x16\xad\xba\xbe\x78\x7f\x1b\x76\xf4\x62\xc1\xc2\x13\x86\x08\x51\xa3\xb1\xe2\xf7\xa0\x46\x72\x6a\xba\x94\x9c\x5c\xaf\x0e\x5a\xf7\xdf\x89\x4e\xf6\x2f\x00\x00\xff\xff\x3c\x4d\xec\x7d\x3b\x00\x00\x00")

func stdlibRandomPiscBytes() ([]byte, error) {
	return bindataRead(
		_stdlibRandomPisc,
		"stdlib/random.pisc",
	)
}

func stdlibRandomPisc() (*asset, error) {
	bytes, err := stdlibRandomPiscBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "stdlib/random.pisc", size: 59, mode: os.FileMode(436), modTime: time.Unix(1491353691, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _stdlibShellPisc = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x54\xcb\x6e\xdb\x3a\x10\x5d\x8b\x5f\x71\xc0\x78\xe1\x2c\x68\x38\x5b\x45\xd7\xc0\x7d\x24\xb7\xd9\xa4\x05\xd2\x9d\x23\x14\x8c\x38\x8e\x88\xd0\x24\x2b\x52\x36\xfa\xf7\xc5\x50\x7e\xb6\x5d\xd8\x20\xcf\xbc\xcf\x19\x4a\xd4\x70\x09\x73\x28\x85\x5b\x38\x9b\xb2\xda\x58\x47\x09\x6b\xc4\xc1\xfa\x8c\x16\x3b\xea\x14\xe9\xae\xc7\xbd\x10\x35\xac\x57\xc6\x0e\x98\x83\xff\xbf\x8f\x21\x73\xe4\x62\x81\x5b\x88\x2a\xee\x0d\xea\x38\xd0\x4e\x54\x6b\x74\x06\x2d\x8c\x8d\xe8\xb4\x73\xa2\x9a\x31\x8e\xce\x88\x92\x25\xf5\xe4\x9c\xda\x06\x33\x3a\xc2\xfc\x94\xe7\x56\x54\xcd\x8e\xba\x1c\x86\x15\xea\xa8\x73\xdf\xdb\x94\x39\xd9\xec\x78\x01\xd7\xe0\x8e\x74\x8c\xe4\xcd\xd9\x0b\x2d\xea\x38\xa6\xde\x5c\xbb\xb3\x6b\x0c\xf1\x4d\x77\x1f\x48\x7b\x1d\x7f\x09\x08\xd1\x08\x51\x35\xc6\x76\x79\x05\x33\x46\xd4\xa9\x07\x67\xb8\x60\xa2\x45\xd3\x28\x97\x18\x15\x55\x75\x45\x91\xa8\xaa\x8a\xa3\xd4\x6a\x1b\x0c\x41\x1a\x89\x94\x07\xd5\x05\x9f\xb5\xf5\x97\x24\xae\x61\x86\x10\xd1\xc2\x6e\x38\xea\x82\x55\x51\x1d\x0a\x98\x52\xb8\xcc\x00\xb9\x58\xc8\x89\xc1\xa6\x51\x63\x3c\x1b\x8e\x58\x57\xe6\xe4\xfe\x4f\x10\xcf\xc8\xe0\x3b\x65\xc5\x5d\xe8\xac\x78\xd6\x83\xbf\x2e\x3c\xca\x57\x7f\x6c\x71\x2c\x7d\x35\x8d\xda\x77\xae\x98\xbe\xfc\xfd\xf5\x93\x04\xf9\x9d\x7a\xa7\x83\x89\x21\xb6\x5d\x81\xe4\x8b\xbe\x6b\xa8\xd9\x9f\xa7\x6e\x27\x3f\x9b\x78\x53\x4a\xea\x17\x96\x1b\xcb\xc5\x9d\x9c\x4c\xdb\x60\xbc\xde\xd2\xb4\x0b\x7a\xff\xa1\xac\xb7\x19\x73\x84\x98\x6d\xf0\xe9\xb8\x52\xa2\xba\xc1\xbf\x03\xe9\x6c\xfd\x3b\x72\x4f\x48\x59\x67\xc2\x9b\x4e\xb6\x4b\xd8\x84\x01\xda\x97\xf0\x94\x7f\x38\xe2\x4a\x2d\xea\x7f\x1e\xfe\x7f\x7a\x3e\x9c\x1f\x9e\xff\x13\x95\x7c\xcd\x12\xf5\xd3\xe3\x0b\xf8\xec\x25\xea\xcf\x8f\x2f\xa2\x92\x12\xf5\x52\x54\x19\xb5\x27\x32\x49\xa5\xe8\x6c\xbe\xdc\x3f\x67\xfd\x31\xe7\x40\x89\xb2\xda\x84\x41\x1d\xc1\xda\x5a\xcc\xf8\xc7\x00\x1c\x79\xac\xfe\xc2\x1a\x52\x16\xad\x27\x94\xed\x65\x55\xf3\xa4\x7b\x3b\xe5\x54\x93\x14\xb3\x8b\xb2\x1c\xb2\xc4\x8c\x7b\x64\x26\x27\xac\x38\xa3\xc5\xbe\x27\xcf\xb1\xe4\xd3\x38\xd0\xb1\x4f\x16\xe5\x7c\x3f\x54\x5c\x9e\xeb\xd5\x77\xbf\xf9\x30\x7a\xda\xf6\x32\x44\xe1\x6a\x92\xe4\x44\x1b\x93\x36\x41\xcc\x1e\x3f\xfa\x1b\x8c\x89\x90\x7b\x9b\x90\x03\x12\xe5\x49\x8b\xa2\xa9\x4e\xe5\xb2\xb7\xb9\x4f\x5d\x88\x24\x6a\x76\x56\x93\xf1\xf0\x55\xb9\x7a\xec\xf2\x1b\xfb\x4a\x4e\xa3\x5c\xe8\xb4\xc3\xfd\xcf\x00\x00\x00\xff\xff\x46\xbd\x3e\x73\x83\x04\x00\x00")

func stdlibShellPiscBytes() ([]byte, error) {
	return bindataRead(
		_stdlibShellPisc,
		"stdlib/shell.pisc",
	)
}

func stdlibShellPisc() (*asset, error) {
	bytes, err := stdlibShellPiscBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "stdlib/shell.pisc", size: 1155, mode: os.FileMode(436), modTime: time.Unix(1504845781, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _stdlibStd_libPisc = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x54\xc1\x6e\xe3\x36\x10\xbd\xfb\x2b\x1e\xb6\x87\x46\x41\x6c\x77\xb3\x37\x1b\xa8\x51\x6c\x2e\x0b\xb4\x48\x8b\x04\xe8\xc1\x35\x1a\x5a\x1e\x45\x84\x69\x52\x25\x47\x56\xfc\xf7\xc5\x0c\x25\x59\x8b\xe4\xb0\x87\xdd\xc8\xe4\xbc\xf7\xe6\x0d\x67\x66\x79\x8b\xdf\x50\x51\x87\x32\x44\x42\x17\xe2\x21\xa1\x0a\x11\x7f\x7e\x7b\xfa\x0a\x7a\xa3\xb2\x65\x1b\xfc\x02\xb7\xcb\xd9\x6c\x85\x26\x5a\xcf\xce\xe3\x06\x06\xf3\x39\x50\xe0\xd7\xc4\xd1\xfa\x57\xb9\x39\xff\xdb\xb4\x9c\xf0\xe9\x1f\xff\x69\xf2\x73\x3d\xe2\x06\xd4\x47\xa0\xf5\x6c\x36\x5b\xde\xe2\xb9\x26\x54\xc1\xb9\xd0\xc9\x75\x69\x9c\x4b\x30\x91\xc0\x35\xa1\x6c\x63\x24\xcf\x48\x6c\xca\xe3\x3c\xd5\x6d\x55\xb9\x9e\xe4\x64\xd9\x9e\x29\x69\x92\xb3\xd5\xc3\xe3\x57\x38\x92\x24\x1d\x79\xb3\x77\x24\xa2\x8e\xfc\x2b\xd7\x28\x80\xdf\xc9\xc3\x26\xa5\x7c\x0d\x38\x50\x65\x3d\x1d\xd4\x39\xb8\x36\x2c\x27\x14\xd3\x00\x28\x6b\x2a\x8f\x22\xc3\x61\x80\xf0\xa5\x21\xa4\x4b\x62\x3a\x41\x13\x57\x45\x7b\x6a\x42\x14\x8f\x95\x75\xd4\x18\xae\x45\x75\xb1\x40\x81\x97\x7c\xf5\x32\x3a\x60\x73\xa4\x04\xa3\x91\x90\xd0\x3b\x18\x7f\x80\x61\xa6\x53\xc3\x49\xa4\x72\xe5\xb3\x71\x0d\x33\x9c\xb3\x8b\xe4\x8c\x98\xbd\x82\x17\x9a\xc3\x4f\x78\x7e\x7c\x78\x5c\xe1\x8f\x70\xce\xa8\xc4\x86\x4f\xe4\x85\x4e\x70\x7d\x7a\x5c\x53\xa2\xde\x8b\x8d\x88\x94\x1a\x2a\x95\xef\x14\x0e\xad\xa3\xa4\xcf\xf0\x24\x25\x46\x2e\xb1\xd4\xe2\x76\x99\x3d\x1e\xda\x66\x78\x44\x03\x83\x02\x0f\x6d\xe3\x6c\x69\x98\x72\x3d\x39\x34\x08\xd5\x20\x5f\x1e\xb1\x9e\xad\xde\x83\x7e\x41\x63\xcb\xe3\x5c\xce\xd7\x7d\xf1\xee\x7f\x84\x99\xbb\x00\x72\x94\x4d\xbd\x93\x19\x29\xf6\x99\x64\xaf\xff\x0a\xdc\x63\x8b\xcf\x57\xc5\x1d\xd8\x9e\x28\x8d\xca\xa9\x33\x13\xd8\x5e\xb5\x9f\x3a\xd3\xfc\xa8\xec\x47\xf8\x41\x2e\x86\xec\x10\xfa\x35\x76\xff\xe0\x9f\x9c\x32\x78\xfb\x1d\xc1\x15\xde\xdf\xdf\x8f\xe8\x7d\xc6\x8b\x23\x3d\xbb\x7a\x59\xe1\xcb\x24\xaa\xcc\x71\x5f\x3e\x8a\x0b\x67\x8a\xb8\xc1\x1b\x2e\x12\x24\x7f\xde\x26\x92\x93\x27\x39\xf4\x69\xfd\xd7\x06\xce\x15\x2d\xf0\xdc\x37\xae\x9c\x19\xd9\x0b\xb9\x6f\x71\x36\xae\xa5\xdc\xc4\x3a\xcd\x52\x22\x3d\xc3\x5e\xea\x14\xfc\xa4\x68\xa6\x62\x8a\x3a\xda\x3a\x55\x35\x4d\xd8\xb4\x58\x47\x22\x51\x5e\x2c\x0c\xde\xf4\x6e\x85\x6d\xff\x6b\xae\x23\xb5\xc7\x0e\xfd\x87\xe4\xae\x8e\xb6\xca\x88\x9d\xa6\x2d\x3e\xf7\x76\xcc\xfe\xb3\xfe\x7f\x9f\x31\x32\x90\xdb\xac\x91\x83\x15\xb7\xd6\xb6\xff\x56\xe9\xd4\xe8\x53\x8f\x3d\x6f\x2b\xdc\x60\x03\x8e\x2d\xfd\x25\xa5\xa8\x8c\x4b\xf9\x6b\x3e\x97\xf9\x69\x1d\xa3\x90\xa8\x8d\x6e\x95\xbb\x4c\x28\xe1\x2a\x7b\x87\xc0\x35\xc5\xce\x26\xca\x37\x23\x7e\x81\x87\x7e\xf7\x98\x84\x97\x8d\xde\xbe\x68\xee\x13\xc9\x1c\xde\x4b\xa1\xc0\x66\xc8\x77\x79\x8b\xbf\x6b\xf2\x93\x84\x65\x6b\x77\x72\x74\x09\x2d\x0e\xc1\xff\xcc\xf0\x24\xec\xce\x49\xdb\xda\x4a\x2d\xe5\x90\x91\x7e\x24\xde\x62\x77\x25\x97\xae\xec\xfb\xc9\xe7\x5e\x7a\xd7\x49\x92\xa6\xd7\xe5\x71\x9d\xdb\x42\x67\x3d\xef\x7a\x8d\x28\x1d\x99\x38\xcf\xef\x7e\x33\x10\xe5\xfd\x2d\x9b\xee\xb2\x81\x0f\x8c\xdd\x95\xbd\xab\x65\xab\xad\x67\xff\x07\x00\x00\xff\xff\x8d\xf7\x0d\x9f\x96\x06\x00\x00")

func stdlibStd_libPiscBytes() ([]byte, error) {
	return bindataRead(
		_stdlibStd_libPisc,
		"stdlib/std_lib.pisc",
	)
}

func stdlibStd_libPisc() (*asset, error) {
	bytes, err := stdlibStd_libPiscBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "stdlib/std_lib.pisc", size: 1686, mode: os.FileMode(436), modTime: time.Unix(1511844549, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _stdlibStringsPisc = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x53\x3d\x6f\xdb\x30\x10\x9d\xa3\x5f\xf1\x1a\x64\x88\x92\x32\xe9\xec\x00\x36\x82\xb4\x43\x26\x0f\xee\x66\x18\x10\x4d\x9d\x23\xa6\x12\x29\x93\x27\xbb\xfe\xf7\xc5\x51\x1f\x56\x3b\x75\xd3\x1d\xdf\xbd\xf7\xee\x43\xcf\x0f\xc8\x36\x1c\xac\xfb\x88\x4f\xad\x8d\x26\x1b\x22\x54\xda\x95\xb5\x7c\x9c\x7d\x28\x23\xb4\x2b\x11\x3b\x53\x65\x0f\xcf\x59\xb6\xf8\xbe\x7e\x43\xe4\xa0\x8c\x77\x46\x33\xee\xa1\xb1\x87\x52\x30\xc8\xf1\x96\x72\xe4\x34\x13\xb8\x22\xb0\x6f\xc1\x67\x2f\x78\x51\x81\x77\x29\x1d\x59\x9b\x5f\x78\xe9\xb9\x96\xfd\x63\x22\x52\x4a\xa0\x3d\xd1\x89\x02\x43\xbb\xcb\x80\x3e\xe9\xba\x23\x58\xc7\x1e\x7a\xe0\x1b\x19\x22\x87\xa5\x75\x62\x45\x8a\x95\x12\xd4\x17\xe4\x78\x65\xa6\xa6\x65\xb0\x87\x19\xf9\xc6\x52\xa1\x71\x02\xa4\x0f\x0a\x33\x22\xf5\xe9\xad\xc3\x3d\x4e\x64\x10\xa9\xbd\x3a\xfa\xd9\x05\x07\x2d\x79\xf6\xe1\x1f\x1f\xfb\x0b\x8a\xa1\x8d\x22\x4d\xab\x8b\x92\x96\x7a\xeb\xb0\x27\x3e\x13\x39\xe8\xba\x4e\xdd\x53\x4d\x0d\x39\x8e\xf0\x07\x14\x27\x32\xc5\x5c\x3e\xb6\xb5\x1d\x3b\x19\xf4\xc5\x4a\x8e\x4d\x7a\x98\x24\x07\x03\x83\x9d\x5e\xaf\x88\xd4\x16\xd0\x11\x1a\x25\xd9\xda\x36\xc4\x7f\xf7\x26\xd3\xb8\xac\xae\x73\x5a\x21\xc7\x7b\x04\x57\x36\x8e\xbc\x03\x64\x5e\x74\x5c\x5d\x97\x2c\x15\xaf\xe1\x3f\x96\x4b\xc7\x4e\xd7\x73\x9e\x65\xe8\x1c\xa9\x40\xba\xa4\x70\x75\xe0\xf7\x9f\xb2\xed\x40\x72\x31\x1a\x73\xcc\x21\xf8\x66\xea\xf7\x2b\x22\xf5\xaa\xc5\xf6\x7d\xbd\x2b\x50\x7a\x13\x71\xf0\x01\x8d\x0f\x72\x16\x07\x3f\x6a\x91\x36\x95\x32\x95\x1e\x55\x8e\x9d\x67\x91\x7a\x7a\x42\x8e\x1f\xbf\xc9\x74\x4c\x28\x24\x5b\x24\x02\xc1\x43\xf0\xda\x30\x85\x7b\xb1\x90\xcb\xda\x8a\xc8\x41\x56\x93\x2d\x86\x73\xef\xa6\x13\x4b\xec\x4a\xc1\x75\xcd\xda\x98\x2e\x44\xe4\x58\x58\xe7\x28\x64\x37\xdf\xb0\xb0\xd9\xcd\x16\x77\x29\x9e\x26\xb8\xc5\xe3\xa3\xc5\x0e\xe7\x8a\x1c\x76\x33\x93\xd9\xcd\x9d\xcd\x46\x15\x99\xb4\x8c\x44\x6e\x38\x6e\x24\xce\xc1\x97\x96\xfc\x01\xb7\xfd\xbf\x79\x3b\x31\x4e\xce\x5c\xbf\x21\xf9\x4c\xe3\x52\x69\x53\x92\xcc\x27\xb0\xf3\x8c\x17\xfc\x09\x00\x00\xff\xff\x74\xb0\x4a\xc5\xf1\x03\x00\x00")

func stdlibStringsPiscBytes() ([]byte, error) {
	return bindataRead(
		_stdlibStringsPisc,
		"stdlib/strings.pisc",
	)
}

func stdlibStringsPisc() (*asset, error) {
	bytes, err := stdlibStringsPiscBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "stdlib/strings.pisc", size: 1009, mode: os.FileMode(436), modTime: time.Unix(1500942096, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _stdlibSymbolsPisc = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\xce\xb1\x0e\x82\x30\x10\xc6\xf1\xbd\x4f\xf1\x1f\x81\x04\xd9\xc5\x48\x0c\xee\x0e\x3e\xc1\xa1\x97\x48\x82\x45\x68\x1b\xe3\xdb\x9b\x4a\x43\x18\x1c\xdb\xfb\x7e\xf7\x5d\x55\x98\xeb\xe7\xd9\x8d\x83\xdb\xbd\x7a\x77\x33\x45\x65\xf6\xe7\x4b\xcb\xc1\xfd\x7e\x8f\x64\x94\x25\x39\xed\xac\xe2\x15\x21\xd8\x7e\x0a\xca\x32\xa6\x5e\xd2\xf1\x55\x5a\x9d\xc8\x10\x3a\x22\x69\xc8\x39\xcd\x8a\x7f\xa8\x53\xfc\x7b\x4c\xc4\x61\x47\x8f\x4e\x41\x06\xa8\xcd\xc6\xaf\x7c\xa3\x05\xb1\x77\xba\x94\x4f\x1b\x9a\x58\xfb\xdf\xac\x87\xc4\x92\xda\x7c\x03\x00\x00\xff\xff\x27\x8e\xbf\xf1\xde\x00\x00\x00")

func stdlibSymbolsPiscBytes() ([]byte, error) {
	return bindataRead(
		_stdlibSymbolsPisc,
		"stdlib/symbols.pisc",
	)
}

func stdlibSymbolsPisc() (*asset, error) {
	bytes, err := stdlibSymbolsPiscBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "stdlib/symbols.pisc", size: 222, mode: os.FileMode(436), modTime: time.Unix(1491353691, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _stdlibVectorsPisc = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x55\x4d\x8b\xe3\x46\x10\x3d\xab\x7f\xc5\x83\x35\x64\xec\x89\x67\xd6\x9a\xdd\x1c\xe4\x10\x08\x31\x84\x9c\x02\x1b\xd8\x8b\x30\x6c\x8f\x55\x8e\x9b\x48\xdd\x4a\x77\x4b\xf6\xfc\xfb\x50\xa5\x4f\xaf\x9d\x4d\x4e\x33\x5d\xf5\xaa\x5e\xd7\xab\xd7\xf2\xf3\x0a\x9f\xe9\x10\x9d\x0f\xd0\xb6\x40\x3c\x91\xf1\x68\xa2\x29\x4d\x34\x14\xb0\x7a\x56\x2a\xdb\xfd\xfe\x0b\x5a\x3a\xac\x03\xc5\xb5\x8e\x78\xe0\x03\x5a\x5d\xc2\x14\x17\xac\xd7\xa0\x92\x2a\x2c\xf1\x07\x45\xce\x44\xe7\xa1\xa3\xe4\xa2\x63\x58\x43\xd8\x4e\x4d\xc6\x06\xd7\xc5\xbf\x52\x94\x7f\xc9\x46\xa9\xb6\x05\x5d\x86\xb2\x1f\xbb\xae\x3f\xe1\x81\x0b\x7a\x8a\x25\x7e\xb3\x26\x42\xc3\xd2\x79\x88\xcd\x68\x48\x1f\x4e\x3d\xd1\xdf\x8d\x8b\x5c\xf8\xf4\x84\x25\x3e\x35\x16\x5f\x38\xf2\x05\x47\xe7\x21\xb0\x81\xd7\x58\x9e\x7f\x6a\xa6\x32\xa4\xfd\xe1\x01\x1a\xaf\x3d\x3b\x96\xd3\x8d\x52\xe4\x08\x67\x5d\x0b\x67\xed\xa9\x26\x5b\x60\x8f\x68\x2a\x0a\x5d\x87\x4d\x18\x06\x9e\x86\x7d\x3f\x28\xb1\x65\x0a\x5b\xdc\x00\x36\x73\xc0\x8b\xbf\x05\xa4\x73\xc0\x87\x78\xba\x01\xbc\xcc\x01\x1f\xef\x00\x3e\xcc\x01\x3f\xdc\x01\x7c\x9c\x00\xef\xb0\x73\xf6\xbb\x08\x5d\x14\xa8\x9c\x27\x9c\xc8\x13\x1a\x5b\x52\x08\x2c\x99\x27\x98\x00\x8d\x55\x4b\xfe\x6d\x85\x3f\x9d\x2b\xe0\x49\x07\x67\x7b\xf7\xb0\xcc\x29\x33\x6c\xd0\xa6\xe3\x42\xda\x17\x2c\xf1\x73\x5d\x97\xec\x34\x2d\x61\x1d\x8d\xb3\xec\x9b\xf9\x62\xce\x26\x10\x6a\x6d\x3c\x6f\xa8\xdd\x88\x51\xdb\xf4\x7b\x86\x79\x0a\x4d\x29\x9b\x6b\x5f\x64\x92\x3b\x4c\xf9\xb0\x3b\xfd\x8a\x7d\x47\xaa\x92\x4c\x52\x59\x9b\x22\x6b\x37\x2a\x79\x8f\xcc\x4c\x5b\x7d\x97\xb9\x26\xaa\x64\xd1\xa6\x28\xc9\x62\xd1\x6e\xe4\x6f\x65\x2c\x72\x95\x24\x7c\x5e\x98\x41\x1e\x09\xa4\x53\x40\x25\x89\xf4\x96\x53\xcd\x86\x50\x49\x92\x63\x83\x47\xec\xb1\xc8\x8c\x4a\xe6\xfe\x60\x75\x42\x5d\xea\x99\x49\x4c\xa4\x2a\x60\x89\x5d\x53\xd5\x62\xc8\x83\xb3\x91\x6c\x0c\x70\xc7\xb9\x41\x9d\x8d\x4e\xce\x21\xea\xc3\x5f\x32\xfd\xbd\x4e\x4f\xe2\xfb\x9c\x47\x1f\xde\xc5\x76\xf6\xa8\x3d\xb5\xe4\x03\x4d\x45\x7d\xa0\x33\xfa\xa7\xee\xc0\xeb\x99\xde\xd8\xff\xa8\x53\x49\xd1\xd4\xc8\xf8\xc0\xca\x65\xfc\x2c\x58\x63\x25\x4a\xac\x45\x09\x51\x66\xc1\x99\x14\xcf\xa2\x4f\xde\x89\x49\x87\x99\xbc\xd9\x65\x8c\x31\x74\x88\xbe\x8d\xd1\xb7\x01\xdc\x7f\x9e\x98\x74\x4c\x5e\xa6\xaa\xeb\xf4\x7a\xdd\x2d\xe6\xf1\x71\x5a\x88\x92\x22\x25\xcf\x96\x2b\x8e\xa6\x8c\xe4\xbf\xfa\x8a\x74\x41\x2a\x3e\xf7\x73\x8e\xa6\xc9\x2c\x9d\xa5\xb7\x58\x8b\x07\x55\x49\x02\x96\xa1\x77\xe1\xa2\x03\x4c\xdf\x8b\xce\x1d\x43\x21\xf6\xc8\x51\x78\x57\x63\x0f\x73\x04\xdf\x6a\x58\x98\x4a\xfa\xda\xe9\x6e\x95\xae\xbf\xba\x58\xc5\xed\xba\x6b\x29\x8c\xfe\xa6\x43\xa7\x3b\x3a\x41\x78\x19\x39\x14\x6e\x64\x16\xf8\x18\x97\x2b\x5e\xcb\xca\x37\x53\x09\x44\x30\x8c\x8a\x61\x94\xec\x79\x85\xd1\x53\x8c\x3d\x7a\x67\x67\x56\xec\x67\x5c\x62\xc7\x13\xb2\x6b\x3b\xc0\xb5\xa5\xb7\xfc\x83\x93\xfd\x67\x0f\xf9\xd8\x0e\x69\x91\xec\xce\x05\x5e\xf9\x51\x7c\x83\x5f\xf2\xdf\xa4\xff\x97\x0e\x3d\xbb\x64\x7b\xf2\x7f\x02\x00\x00\xff\xff\xd5\x00\x21\x41\x46\x07\x00\x00")

func stdlibVectorsPiscBytes() ([]byte, error) {
	return bindataRead(
		_stdlibVectorsPisc,
		"stdlib/vectors.pisc",
	)
}

func stdlibVectorsPisc() (*asset, error) {
	bytes, err := stdlibVectorsPiscBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "stdlib/vectors.pisc", size: 1862, mode: os.FileMode(436), modTime: time.Unix(1512805598, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _stdlibWithPisc = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x90\xcd\x4e\xc3\x30\x10\x84\xcf\xf1\x53\x8c\x4a\x0f\x25\xc2\xcd\x3d\x45\x08\x09\x71\xe0\x82\x10\x57\x84\xd0\xc6\x31\xb1\x8b\x1b\x87\x78\x5b\xbf\x3e\x5a\x97\x22\xfe\x8e\xab\xf1\x7c\x33\xe3\xa6\x46\xf6\xec\xa8\x0b\x16\xb1\xdb\x5a\xc3\x09\x75\xa3\x54\xfb\xf0\x78\x8b\x6b\xac\x90\xe3\xdc\xdf\xd3\xce\x42\x6b\xac\xd7\x38\xc7\xe2\x45\x0c\x0b\x0c\x96\x75\x88\x86\x02\x52\xa6\x09\xbd\x37\xac\x07\xcb\x30\x14\x02\x36\x4a\xb5\x05\x8c\x95\x60\xf1\xbe\x8f\xfc\x45\x68\xcb\xd5\xc6\x6e\xab\xaa\x4b\xf1\x5d\xa1\xcd\x07\x9a\x55\xb5\x2c\xca\x29\x41\x0e\xed\x28\xe9\x03\xcd\x78\xc2\x3f\xa2\x74\x10\xb1\xd8\xf1\x8c\xec\xec\xa8\xaa\xa6\xc6\xdd\x2b\xd8\xf9\x74\xcc\x75\x94\x50\x5c\x48\x96\x2f\x40\xd3\x64\xc7\x1e\x9e\xc1\x11\x84\xc4\x64\xde\x64\xf3\xaf\xf4\xe5\xa9\xb7\x4e\xc7\x14\x55\x9d\xe1\x46\xc6\x95\x67\xfb\xe4\xc7\xe1\x13\xbb\xb3\x4c\xdf\xfa\xcb\x0f\xfc\xa1\x95\x86\x3f\x70\x1b\xf5\x11\x00\x00\xff\xff\x81\x99\x2d\xf3\x7d\x01\x00\x00")

func stdlibWithPiscBytes() ([]byte, error) {
	return bindataRead(
		_stdlibWithPisc,
		"stdlib/with.pisc",
	)
}

func stdlibWithPisc() (*asset, error) {
	bytes, err := stdlibWithPiscBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "stdlib/with.pisc", size: 381, mode: os.FileMode(436), modTime: time.Unix(1491353691, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"stdlib/bools.pisc": stdlibBoolsPisc,
	"stdlib/debug.pisc": stdlibDebugPisc,
	"stdlib/dicts.pisc": stdlibDictsPisc,
	"stdlib/io.pisc": stdlibIoPisc,
	"stdlib/locals.pisc": stdlibLocalsPisc,
	"stdlib/loops.pisc": stdlibLoopsPisc,
	"stdlib/math.pisc": stdlibMathPisc,
	"stdlib/random.pisc": stdlibRandomPisc,
	"stdlib/shell.pisc": stdlibShellPisc,
	"stdlib/std_lib.pisc": stdlibStd_libPisc,
	"stdlib/strings.pisc": stdlibStringsPisc,
	"stdlib/symbols.pisc": stdlibSymbolsPisc,
	"stdlib/vectors.pisc": stdlibVectorsPisc,
	"stdlib/with.pisc": stdlibWithPisc,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"stdlib": &bintree{nil, map[string]*bintree{
		"bools.pisc": &bintree{stdlibBoolsPisc, map[string]*bintree{}},
		"debug.pisc": &bintree{stdlibDebugPisc, map[string]*bintree{}},
		"dicts.pisc": &bintree{stdlibDictsPisc, map[string]*bintree{}},
		"io.pisc": &bintree{stdlibIoPisc, map[string]*bintree{}},
		"locals.pisc": &bintree{stdlibLocalsPisc, map[string]*bintree{}},
		"loops.pisc": &bintree{stdlibLoopsPisc, map[string]*bintree{}},
		"math.pisc": &bintree{stdlibMathPisc, map[string]*bintree{}},
		"random.pisc": &bintree{stdlibRandomPisc, map[string]*bintree{}},
		"shell.pisc": &bintree{stdlibShellPisc, map[string]*bintree{}},
		"std_lib.pisc": &bintree{stdlibStd_libPisc, map[string]*bintree{}},
		"strings.pisc": &bintree{stdlibStringsPisc, map[string]*bintree{}},
		"symbols.pisc": &bintree{stdlibSymbolsPisc, map[string]*bintree{}},
		"vectors.pisc": &bintree{stdlibVectorsPisc, map[string]*bintree{}},
		"with.pisc": &bintree{stdlibWithPisc, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

