package infrastructure

import (
	"os"
	"path"
	"roadmaps/core"
)

type imageManager struct {
	LocalPath string
	UriPath   string
}

func NewImageManager(saveTo string, uriPath string) core.ImageManager {

	if _, err := os.Stat(saveTo); os.IsNotExist(err) {
		os.MkdirAll(saveTo, os.ModePerm)
	}

	return &imageManager{
		LocalPath: saveTo,
		UriPath:   uriPath,
	}
}

func (this *imageManager) Save(data []byte, name string) error {
	p := path.Join(this.LocalPath, name)
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (this *imageManager) GetUrl(name string) string {
	return path.Join(this.UriPath, name)
}
