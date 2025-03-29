package service

import (
	"context"
	"crypto/rand"
	"math/big"
)

type defaultSlugGenerator struct {
	slugLength int
}

func NewSlugGenerator(slugLength int) SlugGenerator {
	return &defaultSlugGenerator{slugLength: slugLength}
}

func (g *defaultSlugGenerator) Generate(_ context.Context) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	b := make([]byte, g.slugLength)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	return string(b), nil
}
