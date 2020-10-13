package domain

import "time"

type User struct {
	Id                string
	Name              string
	NormalizedName    string
	Email             string
	EmailConfirmed    bool
	EmailConfirmation string
	Img               string
	Tokens            []UserToken
	Rights            Permissions
	Pass              []byte
	Salt              []byte
	OAuth             bool
}

type UserToken struct {
	Id          string
	Fingerprint string
	UserAgent   string
	Date        time.Time
}

func (this *User) RemoveToken(i int) {
	this.Tokens[i] = this.Tokens[len(this.Tokens)-1]
	this.Tokens = this.Tokens[:len(this.Tokens)-1]
}

func (this *User) HasRights(r Permissions) bool {
	return Permissions(this.Rights).HasFlag(r)
}
