package common

type LexorankGen interface {
	Between(prevKey, nextKey string) (string, error)
}

type Ranker interface {
	CalculateRank(targetPos int, existingRanks []string) (string, error)
}

type RankerGen struct {
	lexorank LexorankGen
}

func NewRankerGen(l LexorankGen) *RankerGen {
	return &RankerGen{
		lexorank: l,
	}
}

func (r *RankerGen) CalculateRank(targetPos int, existingRanks []string) (string, error) {
	if len(existingRanks) == 0 {
		// Если колонка пуста, создаем первый ранг
		rank, _ := r.lexorank.Between("", "")
		return rank, nil
	}

	// Если вставляем в самое начало (position = 0)
	if targetPos <= 0 {
		rank, _ := r.lexorank.Between("", existingRanks[0])
		return rank, nil
	}

	// Если вставляем в конец (позиция больше или равна количеству задач)
	if targetPos >= len(existingRanks) {
		rank, _ := r.lexorank.Between(existingRanks[len(existingRanks)-1], "")
		return rank, nil
	}

	// Вставляем между двумя задачами
	rank, _ := r.lexorank.Between(existingRanks[targetPos-1], existingRanks[targetPos])
	return rank, nil
}
