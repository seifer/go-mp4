package filter

import (
	"io"
)

type Filter interface {
	// Complex filter
	Filter() error
	// Write mp4 file into writer
	WriteTo(w io.Writer) (int64, error)
}

// EncodeFiltered encodes a media to a writer, filtering the media using the specified filter
func EncodeFiltered(w io.Writer, f Filter) (err error) {
	err = f.Filter()
	if err != nil {
		return err
	}
	_, err = f.WriteTo(w)
	return err
}
