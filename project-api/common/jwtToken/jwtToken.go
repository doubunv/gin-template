package jwtToken

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	RS256 = "RS256"
	HS256 = "HS256"
)

//var privateKey string = "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBdlc3Zk01Z3lYdzcrUVFTMmwyZGFQQ3d1cUVxanNtZWVUUXJRbXVjS09yMFdsU3J6CkZ3bDZqUnR3am1HWEdMVkFIci9TQXYveEFFbmxuUXpLNHVtTURjY24xd2tEeml1ZG9MSzI2eFora0ZWd3A5UXAKTXpKL2xjZnhjdU9ySHlZYzAyUGF2b211TS95ZnJENUN6YW13K1MyYWhqTm1UWTJYbWdCTTBKeGl6YTBNZ0xabQpFbkZHaUI3S09zMjFpTDZBQytVRmhKQUxJSDJMQTBWalpFQ2Q5UnhBcEZJdThnT2V3TnYwSG1Qd0w1ZkN1WEVKCkVpSGhtQ1VZa3phSG04SFJrdi9neWpSbXBDQjdPTFdPNU1PZTlOaGN3bmxUUUZwVGVDakt1eDlVeWpNWC9tMXcKYzV4Y3crV2pveFBVWGd6aXBCY3p6SWZUbmZYWXhBOEVCSitRdlFJREFRQUJBb0lCQUMxYmVyNFQxeHZEcjBWSwpKNTRnOVE0clZoQ2Ryb3h0R3c0L1pLdHlYcFNwNmRwVnFsYjl3Z1kvWjdmdnpnbndmQ1hnc0J4ZnhBalpxTXp0CjN6WnR6VWNQUkk3TWJjalQwSzVGbkY2MXlDQXZGUVphV25NVHlGREl5eHREUUJuNU9qV3pWeEVBVG1ITVB6OHkKdVRRZFNXNmdhTHR4V3pGK2thb2lJRUppNzVWblRGYlJ1T3ZtS1p2WnVRODZwbmhyL1R0TjdyVGFicTNRU2N2QgpPZm9XNWV3R1JUUDNNRTlVSVozVFJRTDBUaEIxeEREbk5zMk9qS0JJeHdYSTRTVFZKOGlKajgranRpanpwby80CjZIeEtUdFp3WklCVHFaWVBTbnVEV2pnU01GVzNtdVJ5aFg1QXZOSG1DWFcyaFN2dzd6UEpNaFFpZ2g4cUhWdXAKMS9NVTlJRUNnWUVBNkdqa3JUZGZsKzkyOFVIQmFOa2Q4ZlFmN2pGbklHS3Rqc0VwR1dIVm1MZ1ZscUxrSEZQKwpUSzltcHowZXRSRW5hSisvYWU0eXBrK2hSSjN1YjlGeElkSUVJOTNXY2pDL1JsWU9mcEorQVRINldNekI0bDFDCjZSbmwxRWZab092dllzZjVkY3UyS0x2bzZhZDFneWE1UzRBNXRZR1cwR3Npdk53MzBuaUNiMFVDZ1lFQTBLazkKMnhhNmFPOWtXNW1sOHp3bUc0TzNuWk9kUnIrb0RoS09BSnQrcVpJdXBXYXpnRFRpeU40bkF3bjRCQ2treGxFYgpKVXRrZGhrZHRwM0dVTnlNcFp5TVo3NlFvU3hKUS9jM2lKR254cGx4TDNsVjZzSHRVUURJNThic2psanlkRGQ1ClhSMWpwQWgyY0pMaDhOdWJjNUJRZ25LSzJZbXRRSHcvQmI3Tmx4a0NnWUVBNDFFOVlpTUlFQTY3a3N2VEJkOWYKQjduVXhXQnkwdDFSanBCbHg1ckdsYUp3RXk5dDJEdGJQdHhNbG9VVWxOdWJaUnppMXhTUlc5UUZFNTA1aXdzYgpFWTVlV0VweHBxMEJXU1Z3OWVJSTl0aFFXaXlqOHVsdk9ab0lEdkxiN1NiM3RQR05rT2JZVzgwM3lkbnQrWUtWCnVFQkpzYnQzckpFdG4vWHhiNTVKVGRVQ2dZQkdYUXZ5MlpNNjE1OXNGRUFWNmU2ZjBLY1dpRFM5S3R3dEFxekkKSDZkeTMwekRrZ0p0OVdlVWZQV0MzTlc1OFhYcE9taUJCdzkxazdxbFhrY0Y5Wm1vTDBaWnBiVjM1RXRBbEJreQpBSXduT0k4bWh6QWwzZ2Q5RFZxeWJpNVBRa2RnVVdTbnRjVm9CMEtKYTc2dFRpRDVzMnl5MEpWcURqZFNTNU9sCkNGYWo0UUtCZ1FEb1hWbURhNVFLRC9oOTlPRGNnWnZXcUdwdWZCN2dTZXpSRDJjSllhNGh1d2N0bWdyR1RRUTYKSm9mYU1ObkNmQUI2aktocFF5L3hMMjkrZkdXUVFZOTFJOFp5UU1RQnR4REVaeGdlYytZaktRS3pDMm43TkVMVgpTb2ZHYnhZQU1JcmxiK2NQeU5WUHp5SVZmLzJ2YXBTbDJCaGRycWE3WFlsbTF4YjFtMXl0N1E9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo="
//var publicKey string = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUF2VzdmTTVneVh3NytRUVMybDJkYQpQQ3d1cUVxanNtZWVUUXJRbXVjS09yMFdsU3J6RndsNmpSdHdqbUdYR0xWQUhyL1NBdi94QUVubG5Reks0dW1NCkRjY24xd2tEeml1ZG9MSzI2eFora0ZWd3A5UXBNekovbGNmeGN1T3JIeVljMDJQYXZvbXVNL3lmckQ1Q3phbXcKK1MyYWhqTm1UWTJYbWdCTTBKeGl6YTBNZ0xabUVuRkdpQjdLT3MyMWlMNkFDK1VGaEpBTElIMkxBMFZqWkVDZAo5UnhBcEZJdThnT2V3TnYwSG1Qd0w1ZkN1WEVKRWlIaG1DVVlremFIbThIUmt2L2d5alJtcENCN09MV081TU9lCjlOaGN3bmxUUUZwVGVDakt1eDlVeWpNWC9tMXdjNXhjdytXam94UFVYZ3ppcEJjenpJZlRuZlhZeEE4RUJKK1EKdlFJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg=="

type (
	Token struct {
		PrivateKey string
		PublicKey  string
		Method     string
	}

	Claims struct {
		TokenClaims *TokenClaims
		jwt.RegisteredClaims
	}

	TokenClaims struct {
		Uid  int64  `json:"uid" bson:"uid"`
		Sign string `json:"sign" bson:"sign"`
	}

	TokenDataResp struct {
		Uid       int64
		Sign      string
		ExpiresAt int64
		IssuedAt  int64
	}

	TokenHandler interface {
		Decrypt(token string) (*TokenDataResp, error) // 验证
		Encrypt(uid int64) (string, string, error)    // 加密

		parseRSAPublicKeyFroPEM() (*rsa.PublicKey, error)
		parseRSAPrivateKeyFroPEM() (*rsa.PrivateKey, error)
		getEncryptMethod() string
	}
)

var NewTokenHandler TokenHandler

func NewToken(privateKey string, publicKey string) TokenHandler {
	return &Token{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Method:     RS256,
	}
}

func (t *Token) getEncryptMethod() string {
	return t.Method
}

func (t *Token) Decrypt(token string) (resp *TokenDataResp, err error) {
	token = strings.Trim(token, "JWT ")
	pubKey, err := t.parseRSAPublicKeyFroPEM()
	if err != nil {
		return
	}

	withClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("令牌校验失败")
		}
		return pubKey, nil
	})
	if err != nil {
		return
	}

	if claims, ok := withClaims.Claims.(*Claims); ok && withClaims.Valid {

		expiresAt := claims.ExpiresAt.Time.Unix()
		issuedAt := claims.IssuedAt.Time.Unix()

		if expiresAt < time.Now().Unix() {
			return nil, errors.New("令牌已失效，请重新登录")
		}

		resp = &TokenDataResp{
			Uid:       claims.TokenClaims.Uid,
			Sign:      claims.TokenClaims.Sign,
			ExpiresAt: expiresAt,
			IssuedAt:  issuedAt,
		}

		return
	}

	return nil, err
}

// Encrypt 加密
func (t *Token) Encrypt(uid int64) (string, string, error) {
	tokenClaims := &TokenClaims{
		Uid:  uid,
		Sign: t.randStr(10),
	}

	pem, err := t.parseRSAPrivateKeyFroPEM()
	if err != nil {
		return "", "", err
	}

	switch t.Method {
	case RS256:
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, &Claims{
			tokenClaims,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 180)), //过期时间
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		})
		signedString, err := token.SignedString(pem)

		return signedString, tokenClaims.Sign, err
	case HS256:
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
			tokenClaims,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		})
		signedString, err := token.SignedString(pem)

		return signedString, tokenClaims.Sign, err
	default:
		return "", "", errors.New("no signature method exists")
	}
}

func (t *Token) parseRSAPublicKeyFroPEM() (*rsa.PublicKey, error) {
	key, err := base64.StdEncoding.DecodeString(t.PublicKey)
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM(key)
}

func (t *Token) parseRSAPrivateKeyFroPEM() (*rsa.PrivateKey, error) {
	key, err := base64.StdEncoding.DecodeString(t.PrivateKey)
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPrivateKeyFromPEM(key)
}

func (t *Token) randStr(strLen int) string {
	letters := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	randBytes := make([]rune, strLen)
	for i := range randBytes {
		randBytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(randBytes)
}
