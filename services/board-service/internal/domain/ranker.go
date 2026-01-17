package domain

type LexorankGen interface {
	Between(prevKey, nextKey string) (string, error)
}
