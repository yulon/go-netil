package netil

import (
	"io"
)

func WriteAll(w io.Writer, b []byte) error {
	max := len(b)
	tn := 0
	for {
		n, err := w.Write(b[tn:])
		if n > 0 {
			tn += n
			if tn == max {
				return nil
			}
		}
		if err != nil {
			return err
		}
	}
}

func ReadFull(r io.Reader, b []byte) error {
	max := len(b)
	tn := 0
	for {
		n, err := r.Read(b[tn:])
		if n > 0 {
			tn += n
			if tn == max {
				return nil
			}
		}
		if err != nil {
			return err
		}
	}
}
