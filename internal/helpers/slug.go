package helpers

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var slugUnsafeChars = regexp.MustCompile(`[^a-z0-9]+`)

func GenerateUniqueSlug(parts ...string) string {
	baseParts := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			baseParts = append(baseParts, trimmed)
		}
	}

	base := strings.ToLower(strings.Join(baseParts, " "))
	base = slugUnsafeChars.ReplaceAllString(base, "-")
	base = strings.Trim(base, "-")
	if base == "" {
		base = "item"
	}

	if len(base) > 180 {
		base = strings.Trim(base[:180], "-")
	}

	id, err := uuid.NewV7()
	suffix := time.Now().UTC().Format("20060102150405")
	if err == nil {
		shortID := strings.Split(id.String(), "-")[0]
		return base + "-" + suffix + "-" + shortID
	}

	return base + "-" + suffix
}
