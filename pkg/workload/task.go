package workload

import "context"

type Task interface {
	Name() string
	CleanUp(ctx context.Context, threadID int) error
	Prepare(ctx context.Context, threadID int) error
	Run(ctx context.Context, threadID int) error
	DBName() string
}

type Options struct {
}
