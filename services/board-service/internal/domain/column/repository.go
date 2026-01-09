package column

import (
	"context"
)

type TaskRepository interface {
	Create(ctx context.Context, task *Column) error
	Update(ctx context.Context, task Column) error
	Delete(ctx context.Context, id ColumnID) error
}
