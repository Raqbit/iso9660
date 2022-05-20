package rockridge

import (
	"github.com/kdomanski/iso9660"
)

const (
	alternateNameEntry = "NM"
)

func GetAlternateName(entries []iso9660.SystemUseEntry) string {
	for _, entry := range entries {
		if entry.Signature == alternateNameEntry {
			return string(entry.Data[1:])
		}
	}

	return ""
}
