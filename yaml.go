package yaml

import (
	"bytes"
	"context"
	"sync"

	"github.com/go-kita/encoding"
	"gopkg.in/yaml.v3"
)

func init() {
	Register(Name)
}

// Name is type name.
const Name = "yaml"

var _ encoding.Marshaler = (*codec)(nil)
var _ encoding.Unmarshaler = (*codec)(nil)

var _bufPool = &sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

type codec struct {
	bufPool *sync.Pool
}

func (c *codec) Marshal(ctx context.Context, v interface{}) ([]byte, error) {
	buf := c.bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		c.bufPool.Put(buf)
	}()
	encoder := yaml.NewEncoder(buf)
	for _, option := range encoderOptionFromContext(ctx) {
		option(encoder)
	}
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	err = encoder.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *codec) Unmarshal(ctx context.Context, data []byte, v interface{}) error {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	for _, option := range decoderOptionFromContext(ctx) {
		option(decoder)
	}
	return decoder.Decode(v)
}

// Register register marshaler/unmarshaler.
func Register(name string) {
	encoding.RegisterMarshaler(name, func() encoding.Marshaler { return &codec{bufPool: _bufPool} })
	encoding.RegisterUnmarshaler(name, func() encoding.Unmarshaler { return &codec{bufPool: _bufPool} })
}
