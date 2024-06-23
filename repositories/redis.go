package repositories

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func (r *Repository) WithRedisTrx(ctx context.Context, keys []string, fn func(pipe redis.Pipeliner) error) error {
	txf := func(tx *redis.Tx) error {
		_, err := tx.Pipelined(ctx, fn)
		return err
	}

	err := r.redis.Watch(ctx, txf, keys...)
	if err == redis.TxFailedErr {
		return r.WithRedisTrx(ctx, keys, fn)
	} else if err != nil {
		return err
	}

	return nil
}
