package local

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
)

type DiskStorage struct {
	folder string
}

func NewDiskStorage(folder string) *DiskStorage {
	return &DiskStorage{folder: folder}
}

func (d *DiskStorage) Write(filename string, data []byte) error {
	err := ioutil.WriteFile(path.Join(d.folder, filename), data, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "error writing file to local disk")
	}
	return nil
}
