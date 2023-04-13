package transformers

import "io"

// Transformer transforms the content read from r and writes the transformed
// content to w.
//
// It MUST return ok as true if the content read from r is changed and written
// to w.
type Transformer func(in io.Reader, out io.Writer) (ok bool, err error)
