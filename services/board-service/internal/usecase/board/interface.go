package board

import (
	"context"
)

type CreateBoard interface {
	Execute(ctx context.Context, req *CreateBoardRequest) (*BoardResponse, error)
}

type GetBoard interface {
	Execute(ctx context.Context, req *GetBoardRequest) (*BoardResponse, error)
}

type UpdateBoard interface {
	Execute(ctx context.Context, req *UpdateBoardRequest) (*BoardResponse, error)
}

type DeleteBoard interface {
	Execute(ctx context.Context, req *DeleteBoardRequest) error
}
