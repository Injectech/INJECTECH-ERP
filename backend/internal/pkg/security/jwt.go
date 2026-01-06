package security

import (
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT claims used in the system.
type Claims struct {
    Subject string   `json:"sub"`
    Roles   []string `json:"roles"`
    jwt.RegisteredClaims
}

// SignToken creates a signed JWT with given secret and ttl.
func SignToken(secret string, subject string, roles []string, ttl time.Duration) (string, time.Time, error) {
    now := time.Now()
    exp := now.Add(ttl)
    claims := Claims{
        Subject: subject,
        Roles:   roles,
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   subject,
            ExpiresAt: jwt.NewNumericDate(exp),
            IssuedAt:  jwt.NewNumericDate(now),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString([]byte(secret))
    if err != nil {
        return "", time.Time{}, err
    }
    return signed, exp, nil
}

// ParseToken validates the token string and returns claims if valid.
func ParseToken(secret, tokenStr string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
        return []byte(secret), nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    return nil, jwt.ErrTokenInvalidClaims
}
