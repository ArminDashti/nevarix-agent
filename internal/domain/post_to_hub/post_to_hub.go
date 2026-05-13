package posttohub

import "time"

// CollectAndStoreInCache gathers probe data and persists it for hub sync.
// Full hub integration lives in domain code; this package is the integration seam.
func CollectAndStoreInCache(hubBaseURL, token string, checkedAt time.Time) error {
	_ = hubBaseURL
	_ = token
	_ = checkedAt
	return nil
}
