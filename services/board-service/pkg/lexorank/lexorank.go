package lexorank

import (
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/morikuni/go-lexorank"
)

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
	// Очищаем входящие ранги от старых солей перед вычислением
	purePrev := strings.Split(prevKey, ":")[0]
	pureNext := strings.Split(nextKey, ":")[0]
	// Вычисляем новый ранг между двумя заданными рангами
	rank, err := l.Gen.Between(lexorank.Key(purePrev), lexorank.Key(pureNext))
	if err != nil {
		return "", err
	}
	// Добавляем новую соль через разделитель
	return fmt.Sprintf("%s:%s", rank, randomString(4)), nil
}

// Генерация случайной строки заданной длины
func randomString(n int) string {
	const letters = "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}
