package model

// A HealthzResponse expresses health check message.
type HealthzResponse struct {
	// フィールドの後ろにいろいろ書いたのがstruct tagっぽい？
	// Rustのderiveマクロみたいなものか？
	Message string `json:"message"`
}
