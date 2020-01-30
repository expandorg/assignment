package externalsvc

import (
	"errors"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	ExpirationKey   = "exp"
	UserIDKey       = "uid"
	IssuerKey       = "iss"
	JWTIDKey        = "jti"
	AudienceKey     = "aud"
	SessionDuration = 8760 * time.Hour // 1 year
)

func parser(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("Unexpected signing method")
	}
	return getJWTSecret(), nil
}

func GenerateSessionJWT(userID uint64) (string, error) {
	expiration := time.Now().Add(SessionDuration).Unix()
	claims := jwt.MapClaims{
		IssuerKey:     os.Getenv("FRONTEND_ADDRESS"),
		ExpirationKey: expiration,
		UserIDKey:     userID,
	}
	return generateJWT(claims, getJWTSecret())
}

func GenerateAPIKeyJWT(userID uint64, jwtID string) (string, error) {
	claims := jwt.MapClaims{
		IssuerKey: os.Getenv("FRONTEND_ADDRESS"),
		JWTIDKey:  jwtID,
		UserIDKey: userID,
	}
	return generateJWT(claims, getJWTSecret())
}

func GenerateWebhookJWT(url string, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		IssuerKey:   os.Getenv("FRONTEND_ADDRESS"),
		AudienceKey: url,
	}
	return generateJWT(claims, secret)
}

func generateJWT(claims jwt.MapClaims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseAPIKeyJWT(tokenString string) (uint64, string, error) {
	claims, err := parseJWT(tokenString)
	if err != nil {
		return 0, "", err
	}

	issuer := claims[IssuerKey].(string)
	if issuer != os.Getenv("FRONTEND_ADDRESS") {
		return 0, "", errors.New("Incorrect issuer")
	}

	jwtID := claims[JWTIDKey].(string)
	userID := uint64(claims[UserIDKey].(float64))
	return userID, jwtID, nil
}

func ParseSessionJWT(tokenString string) (uint64, error) {
	claims, err := parseJWT(tokenString)
	if err != nil {
		return 0, err
	}

	issuer := claims[IssuerKey].(string)
	if issuer != os.Getenv("FRONTEND_ADDRESS") {
		return 0, errors.New("Incorrect issuer")
	}

	expiration := int64(claims[ExpirationKey].(float64))
	if expiration < time.Now().Unix() {
		return 0, errors.New("Authorization token expired")
	}

	userID := uint64(claims[UserIDKey].(float64))
	return userID, nil
}

func parseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, parser)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Unable to parse JWT")
}

func getJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}
