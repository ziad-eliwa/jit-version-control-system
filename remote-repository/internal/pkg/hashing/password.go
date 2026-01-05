package hashing

import "golang.org/x/crypto/bcrypt"

type Password struct {
	Plaintext *string `json:"-"`
	Hash      []byte  `json:"password"`
}

func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)

	if err != nil {
		return err
	}

	p.Plaintext = &plaintextPassword
	p.Hash = hash
	return nil
}

func (p *Password) MatchPassword(password []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, password)

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, err
		} else {
			return false, err
		}
	}
	return true, nil
}
