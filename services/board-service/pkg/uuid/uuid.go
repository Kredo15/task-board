package uuid

import "github.com/google/uuid"

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

// Возвращает обычную строку
func (g *Generator) Generate() string {
	return uuid.New().String()
}
