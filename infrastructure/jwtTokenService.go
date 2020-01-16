package infrastructure

import (
	"fmt"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type JwtTokenService struct {
	UserRepo core.UserRepository
	Secret   string
}

func NewJwtTokenService(ur core.UserRepository, secret string) core.TokenService {
	return &JwtTokenService{ur, secret}
}

func (this *JwtTokenService) Validate(authToken string) (userID string, userName string, rights int, err error) {
	token, err := jwt.ParseWithClaims(authToken, &authClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(this.Secret), nil
	})

	if err != nil {
		return "", "", 0, err
	}

	if claims, ok := token.Claims.(*authClaims); ok && token.Valid {
		if claims.StandardClaims.ExpiresAt < time.Now().Unix() {
			return "", "", 0, core.NewError(core.AuthenticationExpired)
		} else {
			return claims.Id, claims.Name, claims.R, nil
		}
	} else {
		return "", "", 0, core.NewError(core.AuthenticationError)
	}
}

func (this *JwtTokenService) Create(user *domain.User, fingerprint, useragent string) (auth string, refresh string, err error) {
	rid := uuid.New().String()
	auth, err = this.newAuthToken(user, rid, this.Secret)
	refresh, err = this.newRefreshToken(user, rid, this.Secret)

	if user.Tokens == nil {
		user.Tokens = []domain.UserToken{}
	}

	user.Tokens = append(user.Tokens, domain.UserToken{
		Id:          rid,
		Fingerprint: fingerprint,
		UserAgent:   useragent,
		Date:        time.Now()})

	if ok, err := this.UserRepo.Update(user); !ok || err != nil {
		return "", "", err
	}
	return
}

func (this *JwtTokenService) Refresh(authToken, refreshToken, fingerprint, useragent string) (aToken string, rToken string, err error) {

	aClaims, err := this.readAToken(authToken, this.Secret)
	if err != nil {
		return "", "", err
	}

	rClaims, err := this.readRToken(refreshToken, this.Secret)
	if err != nil {
		return "", "", err
	}

	if aClaims.RID != rClaims.RID {
		return "", "", fmt.Errorf("Refresh [%s] and Auth [%s] RID not equals", rClaims.RID, aClaims.RID)
	}

	user := this.UserRepo.Get(rClaims.Id)
	if user == nil {
		return "", "", fmt.Errorf("User not found by ID [%s]", rClaims.Id)
	}

	validRID := false
	validMeta := false
	for i := 0; i < len(user.Tokens); i++ {
		t := user.Tokens[i]
		if t.Id == rClaims.RID && t.Id == aClaims.RID {
			validRID = true
			validMeta = t.Fingerprint == fingerprint && t.UserAgent == useragent
			user.RemoveToken(i)
			break
		}
	}

	if !validMeta || !validRID || len(user.Tokens) >= 10 {
		user.Tokens = user.Tokens[0:0]
		this.UserRepo.Update(user)
		return "", "", fmt.Errorf("Refresh token metadata from client and from db not equals")
	}

	return this.Create(user, fingerprint, useragent)
}

func (this *JwtTokenService) newAuthToken(user *domain.User, rid string, secret string) (token string, err error) {

	claims := &authClaims{
		int(user.Rights),
		user.Id,
		rid,
		user.Name,
		"a",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(168 * time.Hour).Unix(),
			Issuer:    "web",
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

func (this *JwtTokenService) newRefreshToken(user *domain.User, rid string, secret string) (token string, err error) {

	claims := &refreshClaims{
		user.Id,
		rid,
		"r",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3 * 168 * time.Hour).Unix(),
			Issuer:    "web",
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

func (this *JwtTokenService) readRToken(refreshToken string, secret string) (*refreshClaims, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &refreshClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*refreshClaims); ok && token.Valid {
		if claims.StandardClaims.ExpiresAt < time.Now().Unix() {
			return nil, core.NewError(core.AuthenticationExpired)
		} else {
			if claims.Type != "r" {
				return nil, core.NewError(core.AuthenticationError)
			}
			return claims, nil
		}
	} else {
		return nil, core.NewError(core.AuthenticationError)
	}
}

// read without validating because it is already expired
func (this *JwtTokenService) readAToken(aToken string, secret string) (*authClaims, error) {
	token, err := jwt.ParseWithClaims(aToken, &authClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*authClaims); ok {
		if claims.Type != "a" {
			return nil, core.NewError(core.AuthenticationError)
		}
		return claims, nil
	} else {
		return nil, core.NewError(core.AuthenticationError)
	}
}

type authClaims struct {
	R    int    `json:"r"`
	Id   string `json:"id"`
	RID  string `json:"rid"`
	Name string `json:"name"`
	Type string `json:"t"`
	jwt.StandardClaims
}

type refreshClaims struct {
	Id   string `json:"id"`
	RID  string `json:"rid"`
	Type string `json:"t"`
	jwt.StandardClaims
}
