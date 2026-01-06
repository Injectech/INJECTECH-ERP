package security

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes plaintext using bcrypt.
func HashPassword(plain string) (string, error) {
    hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashed), nil
}

// ComparePassword compares plaintext with hashed password.
func ComparePassword(hashed, plain string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
