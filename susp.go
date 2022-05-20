package iso9660

import (
	"bytes"
	"fmt"
	"io"
)

// Very simple SUSP handling, following IEEE P1281.
// This does not handle Continuation Areas (CE) nor root-level SP entries

type SystemUseEntry struct {
	Signature string
	Version   uint8
	Data      []byte
}

func (f *File) SystemUseEntries() ([]SystemUseEntry, error) {
	if f.systemUseEntries != nil {
		return f.systemUseEntries, nil
	}

	entries := make([]SystemUseEntry, 0)

	suReader := bytes.NewReader(f.de.SystemUse)

	// TODO: bytes skipped can be specified in root-level "SP" entry, see 5.3

	for index := 0; ; index++ {
		var n int
		var err error

		sigWord := make([]byte, 2)

		n, err = suReader.Read(sigWord)

		if n != 2 || err == io.EOF {
			return entries, nil
		}

		entryLen, err := suReader.ReadByte()

		if err != nil {
			return entries, fmt.Errorf("could not read system use entry length for entry #%d", index)
		}

		version, err := suReader.ReadByte()

		if err != nil {
			return entries, fmt.Errorf("could not read system use entry version for entry #%d", index)
		}

		// entry length - signature bytes - len byte - version byte
		dataLen := int(entryLen - 4)

		data := make([]byte, dataLen)

		n, err = suReader.Read(data)

		if n != dataLen || err != nil {
			return entries, fmt.Errorf("could not read system use entry data for entry #%d", index)
		}

		entries = append(entries, SystemUseEntry{
			Signature: string(sigWord),
			Version:   version,
			Data:      data,
		})
	}
}
