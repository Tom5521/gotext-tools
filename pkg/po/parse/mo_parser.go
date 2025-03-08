package parse

type MoParser struct {
	data     []byte
	filename string
	errors   []error
}

// func NewMo(path string) (*MoParser, error)
// func NewMoFromReader(r io.Reader, name string) (*MoParser, error)
// func NewMoFromFile(f *os.File) (*MoParser, error)
// func NewMoFromBytes(b []byte, name string) (*MoParser, error)
//
// func (m MoParser) Errors() []error {
// 	return m.errors
// }
//
// func (m MoParser) Parse() *po.File
