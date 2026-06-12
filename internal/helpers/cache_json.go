package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	cacheDomain "saetechnology-be/internal/domain/cache"
)

const DefaultCacheTTL = 5 * time.Minute

type PaginatedCache[T any] struct {
	Data  []T   `json:"data"`
	Total int64 `json:"total"`
}

func PaginationCacheKey(
	prefix string,
	page int,
	limit int,
	search string,
	searchBy string,
	orderBy string,
	orderType string,
) string {
	return fmt.Sprintf(
		"%s:list:page=%d:limit=%d:search=%s:search_by=%s:order_by=%s:order_type=%s",
		prefix,
		page,
		limit,
		search,
		searchBy,
		orderBy,
		orderType,
	)
}

func GetJSONCache[T any](
	ctx context.Context,
	repository cacheDomain.Repository,
	key string,
) (*T, bool) {
	value, err := repository.Get(ctx, key)
	if err != nil {
		return nil, false
	}

	var result T

	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return nil, false
	}

	return &result, true
}

func SetJSONCache(
	ctx context.Context,
	repository cacheDomain.Repository,
	key string,
	value interface{},
	ttl time.Duration,
) {
	if value == nil {
		return
	}

	reflected := reflect.ValueOf(value)
	if reflected.Kind() == reflect.Ptr && reflected.IsNil() {
		return
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return
	}

	_ = repository.Set(
		ctx,
		key,
		string(bytes),
		ttl,
	)
}

func DeleteCacheByPrefix(
	ctx context.Context,
	repository cacheDomain.Repository,
	prefix string,
) {
	_ = repository.DeleteByPrefix(ctx, prefix)
}
