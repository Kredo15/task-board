package board

import (
	"context"
)

type TaskRank struct {
	ID   TaskID
	Rank string
}

type ColumnRank struct {
	ID   ColumnID
	Rank string
}

type BoardRepository interface {
	// Create сохраняет новую доску.
	Create(ctx context.Context, b *Board, event *DomainEvent) error

	// Update обновляет метаданные доски (title, description).
	Update(ctx context.Context, b *Board, event *DomainEvent) error

	// Delete удаляет доску и все связанные с ней колонки и задачи.
	Delete(ctx context.Context, id BoardID, event *DomainEvent) error

	// GetByID получает "легкую" версию доски (без колонок).
	GetByID(ctx context.Context, id BoardID) (*Board, error)

	// GetFullBoard получает всё дерево (Board -> Columns -> Tasks) через JOIN.
	GetFullBoard(ctx context.Context, id BoardID) (*Board, error)

	// GetColumnByID получает колонку по ID.
	GetColumnByID(ctx context.Context, boardID BoardID, columnID ColumnID) (*Column, error)

	// SaveColumn создает или обновляет колонку.
	SaveColumn(ctx context.Context, boardID BoardID, c *Column, event *DomainEvent) error

	// DeleteColumn удаляет колонку.
	DeleteColumn(ctx context.Context, boardID BoardID, colID ColumnID, event *DomainEvent) error

	// SaveTask создает или обновляет задачу.
	SaveTask(ctx context.Context, boardID BoardID, t *Task, event *DomainEvent) error

	// GetTaskByID получает задачу по ID.
	GetTaskByID(ctx context.Context, boardID BoardID, taskID TaskID) (*Task, error)

	// DeleteTask удаляет задачу.
	DeleteTask(ctx context.Context, boardID BoardID, taskID TaskID, event *DomainEvent) error

	// GetColumnRanks возвращает список всех существующих рангов колонок на доске.
	GetColumnRanks(ctx context.Context, boardID BoardID) ([]ColumnRank, error)

	// GetTaskRanks возвращает список рангов всех задач в конкретной колонке.
	GetTaskRanks(ctx context.Context, bID BoardID, colID ColumnID) ([]TaskRank, error)
}
