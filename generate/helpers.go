package generate

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cscoding21/csmig/shared"
)

func getMigrationName() string {
	timestamp := fmt.Sprintf("m%s", strconv.FormatInt(time.Now().UTC().UnixNano(), 10))

	return timestamp
}

func NewMigrationObject(description string) shared.Migration {
	return shared.Migration{
		Name:        getMigrationName(),
		Description: description,
	}
}
