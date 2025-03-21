package utility

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

// Userclaim struct defines the that jwt token holds
type UserClaim struct {
	UserID      uint
	Username    string
	Role        string
	PayloadHash string
	jwt.StandardClaims
}

// GenerateToken will generate token for 5 hours with given data
func GenerateToken(key, username string, userID uint) (string, error) {
	expTime := time.Now().Add(time.Hour * 5).Unix()

	claims := &UserClaim{
		UserID:      userID,
		Username:    username,
		Role:        "user",
		PayloadHash: hashPayload(username, userID),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime,
			Subject:   username,
			IssuedAt:  time.Now().Unix(),
		},
	}
	fmt.Println("Generated PayloadHash:", hashPayload(username, userID))

	//Use the combination of email and key for signing
	siginKey := []byte(key)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString(siginKey)
	if err != nil {
		log.Printf("unable to generate token for user %v, err: %v", username, err.Error())
		return "", err
	}
	return signedToken, nil
}

func hashPayload(email string, userID uint) string {
	h := sha256.New()
	h.Write([]byte(email + fmt.Sprint(userID)))
	return hex.EncodeToString(h.Sum(nil))
}
