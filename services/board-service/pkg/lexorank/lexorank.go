package lexorank

import "github.com/morikuni/go-lexorank"

type LexorankGen struct {
	Gen *lexorank.Generator
}

func NewLexorankGen() *LexorankGen {
	generator := lexorank.NewGenerator()
	return &LexorankGen{
		Gen: generator,
	}
}

func (l *LexorankGen) Between(prevKey, nextKey string) (string, error) {
	return l.Between(prevKey, nextKey)
}
