package yaml

import (
	"context"
	"io"
	"math"
	"reflect"
	"testing"

	"github.com/go-kita/encoding"
)

func Test_codec_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		v       interface{}
		want    []byte
		wantErr bool
	}{
		{
			name:    "nil",
			v:       nil,
			want:    []byte("null\n"),
			wantErr: false,
		},
		{
			name:    "empty_struct",
			v:       &struct{}{},
			want:    []byte("{}\n"),
			wantErr: false,
		},
		{
			name:    "map_str_str",
			v:       map[string]string{"v": "hi"},
			want:    []byte("v: hi\n"),
			wantErr: false,
		},
		{
			name:    "map_str_any",
			v:       map[string]interface{}{"v": "hi"},
			want:    []byte("v: hi\n"),
			wantErr: false,
		},
		{
			name:    "map_str_bool",
			v:       map[string]bool{"v": true},
			want:    []byte("v: true\n"),
			wantErr: false,
		},
		{
			name:    "map_str_int32",
			v:       map[string]int32{"v": 10},
			want:    []byte("v: 10\n"),
			wantErr: false,
		},
		{
			name:    "map_str_int64",
			v:       map[string]int64{"v": math.MaxUint32 + 1},
			want:    []byte("v: 4294967296\n"),
			wantErr: false,
		},
		{
			name:    "map_str_int64_large",
			v:       map[string]int64{"v": math.MaxInt64},
			want:    []byte("v: 9223372036854775807\n"),
			wantErr: false,
		},
		{
			name:    "map_str_float",
			v:       map[string]float64{"v": 0.1},
			want:    []byte("v: 0.1\n"),
			wantErr: false,
		},
		{
			name:    "map_str_float_inf",
			v:       map[string]float64{"v": math.Inf(+1)},
			want:    []byte("v: .inf\n"),
			wantErr: false,
		},
		{
			name:    "map_str_float_neg_inf",
			v:       map[string]float64{"v": math.Inf(-1)},
			want:    []byte("v: -.inf\n"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &codec{
				bufPool: _bufPool,
			}
			got, err := c.Marshal(context.Background(), tt.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("codec.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("codec.Marshal() = %v(%s), want %v(%s)", got, got, tt.want, tt.want)
			}
		})
	}
}

func Test_codec_Unmarshal(t *testing.T) {
	type args struct {
		data []byte
		v    interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "blank",
			args: args{
				data: []byte(``),
				v:    (*struct{})(nil),
			},
			wantErr: false,
		},
		{
			name: "empty_struct",
			args: args{
				data: []byte(`{}`),
				v:    &struct{}{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &codec{
				bufPool: _bufPool,
			}
			if err := c.Unmarshal(context.Background(), tt.args.data, tt.args.v); (err != nil && err != io.EOF) != tt.wantErr {
				t.Errorf("codec.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "yml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Register(tt.name)
			if encoding.GetMarshaler(tt.name) == nil {
				t.Errorf("can not get Marshaler after register %s", tt.name)
			}
			if encoding.GetUnmarshaler(tt.name) == nil {
				t.Errorf("can not get Unmarshaler after register %s", tt.name)
			}
		})
	}
}
