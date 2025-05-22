package engine

import "context"

type ctxkey int

const (
	QuotaInventoryKey ctxkey = iota
)

func ContextSetQuota(ctx context.Context, quotas QuotaInventory) context.Context {
	// Check if the context is nil
	if ctx == nil {
		return context.Background()
	}

	// Set the QuotaInventory in the context
	return context.WithValue(ctx, QuotaInventoryKey, &quotas)
}

func ContextGetQuotaByName(ctx context.Context, name string) (Quota, bool) {
	if name == "" {
		return nil, true // No quota name provided, this is not an error
	}

	// Check if the context is nil
	if ctx == nil {
		return nil, false
	}

	// Retrieve the QuotaInventory from the context
	quotas, ok := ctx.Value(QuotaInventoryKey).(*QuotaInventory)
	if !ok {
		return nil, false
	}

	// Get the quota by name
	quota, exists := (*quotas)[name]
	return quota, exists
}
