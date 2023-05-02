package token

import (
	"log"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

type TokenMaker interface {
	CreateToken(memberID int32, duration time.Duration) (string, *paseto.Token, error)
	VerifyToken(signedToken string) (*paseto.Token, error)
}

type PasetoTokenMaker struct {
	Parser    *paseto.Parser
	SecretKey paseto.V4SymmetricKey
}

func NewTokenMaker(key string) (TokenMaker, error) {
	secretKey, err := paseto.V4SymmetricKeyFromBytes([]byte(key))
	if err != nil {
		log.Print("failded to generate secret key: ", err)
		return nil, err
	}

	parser := paseto.NewParser()

	return &PasetoTokenMaker{
		Parser:    &parser,
		SecretKey: secretKey,
	}, nil
}

func (maker *PasetoTokenMaker) CreateToken(memberID int32, duration time.Duration) (string, *paseto.Token, error) {
	token := paseto.NewToken()
	now := time.Now().UTC().Truncate(time.Second)
	exp := now.Add(duration).Truncate(time.Second)

	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(exp)
	token.SetJti(uuid.New().String())
	token.Set("member_id", memberID)

	return token.V4Encrypt(maker.SecretKey, nil), &token, nil
}

func (maker *PasetoTokenMaker) VerifyToken(encrypted string) (*paseto.Token, error) {
	maker.Parser.AddRule(paseto.NotExpired())
	maker.Parser.AddRule(paseto.ValidAt(time.Now().UTC().Truncate(time.Second)))

	parsedToken, err := maker.Parser.ParseV4Local(maker.SecretKey, encrypted, nil)
	if err != nil {
		return nil, err
	}

	return parsedToken, nil
}
