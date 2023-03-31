package transformers

import "io"

// Transformer transforms the content of the file. It MUST return ok as true if
// the content read from r is modified and written to w.
type Transformer func(in io.Reader, out io.Writer) (ok bool, err error)
