package pkg

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JWTCliams struct {
	UserID uuid.UUID `json:"user_id"`
	Tier   string    `json:"tier"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uuid.UUID, tier string, secretKey []byte) (string, error) {
	claims := JWTCliams{
		UserID: userID,
		Tier:   tier,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(jwt.TimeFunc().Add(24 * 7 * 60 * 60)),
			IssuedAt:  jwt.NewNumericDate(jwt.TimeFunc()),
			Issuer:    "GOD Diaku",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := []byte(secretKey)
	return token.SignedString(secret)
}

func ValidateToken(tokenString string, secretKey []byte) (*JWTCliams, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCliams{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTCliams); ok && token.Valid {
		return claims, nil
	} else {
		return nil, jwt.ErrSignatureInvalid
	}
}
