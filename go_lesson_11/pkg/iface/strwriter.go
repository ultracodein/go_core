package iface

import (
	"io"
)

// WriteString пишет во writer только строки
func WriteString(writer *io.Writer, items ...interface{}) error {
	for _, item := range items {
		str, isString := item.(string)
		if !isString {
			continue
		}

		_, err := (*writer).Write([]byte(str))
		if err != nil {
			return err
		}
	}
	return nil
}
