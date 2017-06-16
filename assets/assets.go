// Code generated by go-bindata.
// sources:
// static/css/style.css
// static/index.html
// static/ip.html
// DO NOT EDIT!

package assets

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

var _cssStyleCss = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\xcc\x41\x0a\x83\x30\x10\x05\xd0\x7d\x4e\xf1\xa1\xeb\x40\xbb\x70\x13\x41\xe8\x51\xa2\xc6\x38\x34\xce\x84\xe9\xb4\x28\xc5\xbb\x17\x5d\x75\xd1\x03\xbc\xd7\xcb\xb8\xe1\xe3\x80\x3e\x0e\x8f\xac\xf2\xe2\xd1\x0f\x52\x44\x03\x88\xe7\xa4\x64\xad\x03\x26\x61\xf3\x53\x5c\xa8\x6c\x01\x77\xa5\x58\x5a\xb7\x3b\x77\xe2\x0e\x23\xbd\x2f\x54\xcf\x65\x89\x9a\x89\xbd\x49\x0d\x68\xae\x75\x3d\xb0\xa5\xd5\x7c\x2c\x94\x39\x60\x48\x6c\x49\xff\xe0\x0e\x15\xcf\x1a\xf9\x27\x51\xca\xb3\x05\xdc\x9a\xa3\xd9\xdd\x37\x00\x00\xff\xff\x96\xd2\xd8\x95\xaa\x00\x00\x00")

func cssStyleCssBytes() ([]byte, error) {
	return bindataRead(
		_cssStyleCss,
		"css/style.css",
	)
}

func cssStyleCss() (*asset, error) {
	bytes, err := cssStyleCssBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "css/style.css", size: 170, mode: os.FileMode(420), modTime: time.Unix(1497618037, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _indexHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func indexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_indexHtml,
		"index.html",
	)
}

func indexHtml() (*asset, error) {
	bytes, err := indexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "index.html", size: 0, mode: os.FileMode(420), modTime: time.Unix(1497468196, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _ipHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x92\x4f\x6b\xdc\x30\x10\xc5\xef\xfd\x14\x13\x5d\x7a\xa9\xd6\x14\x0a\xd9\x06\xc9\x25\xa4\x7b\xf0\xa9\x86\xa6\x7f\x42\xc8\x41\x91\xc6\xd1\xb4\xf2\x9f\x4a\xb3\x5e\x96\x25\xdf\xbd\xc4\xce\xc6\xde\x6e\x52\xe2\x8b\xa4\xf7\xe6\xfd\x1e\x0c\x56\x27\x9f\xbf\x5c\x5c\x5e\x95\x2b\xf0\x5c\x87\xfc\x8d\x3a\x91\xf2\x9a\x2a\x08\x0c\xc5\x0a\x4e\x6f\x72\x18\x3e\xf5\xe0\x82\x0d\x26\x25\x2d\x9a\x56\xfe\x4a\x10\x58\x12\x7e\x1c\x8f\xe5\x78\x9c\x8a\x1c\xd4\xc9\x35\x36\x8e\xaa\x1b\x29\x27\xda\x1c\xf5\x0a\xda\x7f\x30\xcb\xd7\x60\x5e\xca\xdf\xf1\x01\xe2\x38\x7f\x14\x1c\x97\x32\x0e\xa3\x71\xe3\x15\x40\xd5\xc8\x06\xac\x37\x31\x21\x6b\xb1\xe6\x4a\x2e\xc5\xa1\xe9\x99\x3b\x89\x7f\xd6\xd4\x6b\xf1\x53\x7e\x3b\x97\x17\x6d\xdd\x19\xa6\xdb\x80\x02\x6c\xdb\x30\x36\xac\x45\xb1\xd2\xe8\xee\x70\xca\x32\x71\xc0\xfc\x87\x37\xfc\x36\x41\xbd\x85\xa2\xfc\xa4\xb2\x51\x3c\xc0\x37\xa6\x46\x2d\x7a\xc2\x4d\xd7\x46\x9e\x11\x37\xe4\xd8\x6b\x87\x3d\x59\x94\xc3\xe3\x1d\x50\x43\x4c\x26\xc8\x64\x4d\x40\xfd\x7e\x6a\x0b\xd4\xfc\x86\x88\x41\x8b\xc4\xdb\x80\xc9\x23\xb2\x00\x1f\xb1\x7a\x50\x0c\x93\xcd\x6c\x4a\xd9\x60\x2e\x6c\x4a\x8f\x49\x95\x4d\xcb\x50\xb7\xad\xdb\x3e\x5e\x1d\xf5\x40\x4e\x0b\xea\x9e\x2a\x00\xae\xda\x75\x84\xa2\x04\x4a\x67\x4f\x9a\xea\x26\x7f\x78\xa7\xce\x34\x79\x51\x82\x71\x2e\x62\x4a\x67\x2a\x1b\x94\xdd\x0e\x16\x45\x79\x3e\x8a\xdf\x3f\xc0\xfd\xfd\x84\xc8\x66\x8c\xe7\x79\x97\x1e\xf7\x40\xd8\x98\x04\x15\xb2\xf5\xe8\xa0\x8a\x6d\x0d\xec\x11\x52\xbb\x8e\x16\xe7\x65\x5f\x07\xe5\xc5\x9e\xdd\x0e\xa8\x82\xc5\x2a\xc6\x36\x1e\x0c\x3d\xdb\x3f\x8c\xcd\xe9\xc7\xb9\x7f\xe0\xd8\xb8\xbd\xad\x32\x47\xfd\x7e\xdd\xe3\x8e\x55\x36\xfe\x8e\x7f\x03\x00\x00\xff\xff\x32\xd9\x5a\xb8\xb4\x03\x00\x00")

func ipHtmlBytes() ([]byte, error) {
	return bindataRead(
		_ipHtml,
		"ip.html",
	)
}

func ipHtml() (*asset, error) {
	bytes, err := ipHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "ip.html", size: 948, mode: os.FileMode(420), modTime: time.Unix(1497617750, 0)}
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
	"css/style.css": cssStyleCss,
	"index.html": indexHtml,
	"ip.html": ipHtml,
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
	"css": &bintree{nil, map[string]*bintree{
		"style.css": &bintree{cssStyleCss, map[string]*bintree{}},
	}},
	"index.html": &bintree{indexHtml, map[string]*bintree{}},
	"ip.html": &bintree{ipHtml, map[string]*bintree{}},
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
