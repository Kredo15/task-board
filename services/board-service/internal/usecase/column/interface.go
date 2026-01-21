package column

import (
	"context"
)

type AddColumn interface {
	Execute(ctx context.Context, req *AddColumnRequest) (*ColumnResponse, error)
}

type MoveColumn interface {
	Execute(ctx context.Context, req *MoveColumnRequest) (*ColumnResponse, error)
}

type UpdateColumn interface {
	Execute(ctx context.Context, req *UpdateColumnRequest) (*ColumnResponse, error)
}

type DeleteColumn interface {
	Execute(ctx context.Context, req *DeleteColumnRequest) error
}
