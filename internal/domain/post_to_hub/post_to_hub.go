package posttohub

import "time"

// CollectAndStoreInCache pulls metrics and forwards them to the hub (placeholder).
func CollectAndStoreInCache(hubBaseURL, token string, checkedAt time.Time) error {
	_ = hubBaseURL
	_ = token
	_ = checkedAt
	return nil
}
