package cli

import (
	"fmt"
	"io"
)

type writemeta struct {
	w   io.Writer
	txt string
}

type lazyWriter struct {
	stdout, stderr io.Writer
	input          []writemeta
}

func (lw *lazyWriter) flush() {
	for _, line := range lw.input {
		fmt.Fprint(line.w, line.txt)
	}
}

type stdoutWriter lazyWriter

func (w *stdoutWriter) Write(b []byte) (int, error) {
	l := writemeta{w.stdout, string(b)}
	w.input = append(w.input, l)
	return len(b), nil
}

type stderrWriter lazyWriter

func (w *stderrWriter) Write(b []byte) (int, error) {
	l := writemeta{w.stderr, string(b)}
	w.input = append(w.input, l)
	return len(b), nil
}
