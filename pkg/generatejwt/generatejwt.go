package generatejwt

import (
	config "djson/internal/config"
	"djson/internal/core/session"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// load jwtKey, from config
var jwtKey []byte

func init() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("jwt key not found")
	}
	jwtKey = []byte(cfg.JWT_KEY)
}

// var jwtKey = []byte("ajfdalkjfalkadfas")

type Claims struct {
	Uid int `json:"uid"`
	jwt.StandardClaims
}

func GenerateJWT(uid int) (string, time.Time, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", time.Now(), err
	}
	return tokenString, expirationTime, nil
}

func GenerateSession(userId int) (session.Session, error) {
	tokenId := uuid.New()
	expiresAt := time.Now().Add(2 * time.Hour)
	issuedAt := time.Now()
	hashToken, err := bcrypt.GenerateFromPassword([]byte(tokenId.String()), bcrypt.DefaultCost)
	if err != nil {
		return session.Session{}, err
	}

	session := session.Session{
		Id:        tokenId,
		Uid:       userId,
		TokenHash: string(hashToken),
		ExpiresAt: expiresAt,
		IssuedAt:  issuedAt,
	}
	return session, nil
}

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, err
		}
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return claims, nil
}
