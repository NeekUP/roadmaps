package infrastructure

import (
	"bytes"
	"github.com/NeekUP/roadmaps/core"
	"os"
	"path"

	"github.com/nullrocks/identicon"
)

type imageManager struct {
	LocalPath       string
	UriPath         string
	AvatarPath      string
	AvatarGenerator *identicon.Generator
}

func NewImageManager(saveTo string, uriPath string) core.ImageManager {

	if _, err := os.Stat(saveTo); os.IsNotExist(err) {
		err := os.MkdirAll(saveTo, os.ModePerm)
		if err != nil {
			panic("fail to create image folder ")
		}
	}

	if _, err := os.Stat(path.Join(saveTo, "users")); os.IsNotExist(err) {
		err := os.MkdirAll(path.Join(saveTo, "users"), os.ModePerm)
		if err != nil {
			panic("fail to create user avatar folder ")
		}
	}

	generator, _ := identicon.New(
		"avatar", // Namespace
		7,        // Number of blocks (Size)
		5,        // Density
	)
	generator.Option()
	return &imageManager{
		LocalPath:       saveTo,
		UriPath:         uriPath,
		AvatarPath:      "users",
		AvatarGenerator: generator,
	}
}

func (mananger *imageManager) GenerateAvatar(username string) ([]byte, error) {
	ii, err := mananger.AvatarGenerator.Draw(username)
	if err != nil {
		return nil, err
	}
	ii.FillColor.RGBA()
	var buf bytes.Buffer
	err = ii.Png(160, &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (mananger *imageManager) SaveAvatar(data []byte, name string) error {
	dir := path.Join(mananger.LocalPath, mananger.AvatarPath, name)
	f, err := os.Create(dir)
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

func (mananger *imageManager) GetAvatarUrl(name string) string {
	if name == "" {
		return ""
	}
	return path.Join(mananger.UriPath, mananger.AvatarPath, name)
}

func (mananger *imageManager) SaveResourceCover(data []byte, name string) error {
	p := path.Join(mananger.LocalPath, name)
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

func (mananger *imageManager) GetResourceCoverUrl(name string) string {
	if name == "" {
		return ""
	}
	return path.Join(mananger.UriPath, name)
}
