package auth

import (
	"dubhe-ci/errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	TokenValid = errors.New(999901)
	SignKey    string
	Expired    int
)

//JWT签名结构
type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{SigningKey: []byte(GetSignKey())}
}

func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

func GetSignKey() string {
	return SignKey
}

func SetExpired(expired int) {
	Expired = expired
}

type Claims struct {
	UserId int `json:"userId"`
	RoleId int `json:"roleId"`
	jwt.StandardClaims
}

func (j *JWT) CreateToken(claims Claims) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(Expired) * time.Second).Unix()
	claims.ExpiresAt = expiresAt
	claims.IssuedAt = now.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenValid
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenValid
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenValid
			} else {
				return nil, TokenValid
			}
		}
	}

	if token == nil {
		return nil, TokenValid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenValid
}

//更新Token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}

	return "", TokenValid
}
