package migrate

import (
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"crypto/sha256"
)

func extractVersionDescription(name string) (uint64, string, error) {
	base := filepath.Base(name)

	if ext := filepath.Ext(base); ext != ".go" {
		return 0, "", fmt.Errorf("can not extract version from %q", base)
	}

	idx := strings.IndexByte(base, '_')
	if idx == -1 {
		return 0, "", fmt.Errorf("can not extract version from %q", base)
	}

	version, err := strconv.ParseUint(base[:idx], 10, 64)
	if err != nil {
		return 0, "", err
	}

	description := base[idx+1 : len(base)-len(".go")]

	return version, description, nil
}

func calculateHash(migrationSource *string) (string, error) {
	if migrationSource == nil {
		return "", fmt.Errorf("can not calculate hash from nil source")
	}
	hash := sha256.New()
	hash.Write([]byte(*migrationSource))
	return hex.EncodeToString(hash.Sum(nil)), nil
}
