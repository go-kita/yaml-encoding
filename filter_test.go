package yaml

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-kita/encoding"
)

func TestWithDecoderOption(t *testing.T) {
	type args struct {
		unmarshaler encoding.Unmarshaler
		opt         []DecoderOption
	}
	tests := []struct {
		name string
		args args
		want func(encoding.Unmarshaler) bool
	}{
		{
			name: "knownFields",
			args: args{
				unmarshaler: &codec{bufPool: _bufPool},
				opt:         []DecoderOption{OnlyAllowKnownFields(true)},
			},
			want: func(unmarshaler encoding.Unmarshaler) bool {
				type N struct {
					F string `yaml:"f"`
				}
				err := unmarshaler.Unmarshal(context.Background(), []byte("k: true"), &N{})
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDecoderOption(tt.args.unmarshaler, tt.args.opt...); !tt.want(got) {
				t.Fail()
			}
		})
	}
}

func TestWithEncoderOption(t *testing.T) {
	type args struct {
		marshaler encoding.Marshaler
		opt       []EncoderOption
	}
	tests := []struct {
		name string
		args args
		want func(encoding.Marshaler) bool
	}{
		{
			name: "indent",
			args: args{
				marshaler: &codec{bufPool: _bufPool},
				opt:       []EncoderOption{SetIndent(4)},
			},
			want: func(marshaler encoding.Marshaler) bool {
				type In struct {
					B bool `yaml:"b"`
				}
				type Out struct {
					I In `yaml:"i"`
				}
				bytes, err := marshaler.Marshal(context.Background(), &Out{I: In{B: false}})
				if err != nil {
					t.Logf("expect not error, got %v", err)
					return false
				}
				want := []byte("i:\n    b: false\n")
				if !reflect.DeepEqual(bytes, want) {
					t.Logf("expect %v(%s), got %v(%s)", want, want, bytes, bytes)
					return false
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithEncoderOption(tt.args.marshaler, tt.args.opt...); !tt.want(got) {
				t.Fail()
			}
		})
	}
}
