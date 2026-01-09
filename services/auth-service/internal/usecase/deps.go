package usecase

// Интерфейсы для зависимостей, которые реализует Infrastructure
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type TokenManager interface {
	GeneratePair(userID string) (access, refresh string, err error)
}
