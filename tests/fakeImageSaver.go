package tests

type fakeImageManager struct{}

func (this *fakeImageManager) SaveResourceCover(data []byte, name string) error {
	return nil
}

func (this *fakeImageManager) GetUrl(name string) string {
	return name
}
