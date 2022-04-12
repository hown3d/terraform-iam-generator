package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type obj struct {
	key string
}

func (o obj) String() string {
	return o.key
}

var objFromStringFunc = func(s string) (obj, error) {
	return obj{key: s}, nil
}

func TestSliceFlag_String(t *testing.T) {
	tests := []struct {
		name   string
		values []obj
		want   string
	}{
		{
			name:   "1 value",
			values: []obj{{key: "hello"}},
			want:   "hello",
		},

		{
			name:   "multiple values",
			values: []obj{{key: "hello"}, {key: "world"}},
			want:   "hello, world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSliceFlag(objFromStringFunc)
			s.values = tt.values
			stringRepr := s.String()
			assert.Equal(t, tt.want, stringRepr)
		})
	}
}

func TestSliceFlag_Set(t *testing.T) {
	tests := []struct {
		name    string
		values  []string
		want    []obj
		wantErr bool
	}{
		{
			name:   "1 value",
			values: []string{"hello"},
			want: []obj{
				{
					key: "hello",
				},
			},
			wantErr: false,
		},
		{
			name:   "multiple values",
			values: []string{"hello", "world"},
			want: []obj{
				{
					key: "hello",
				},

				{
					key: "world",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSliceFlag(objFromStringFunc)
			for _, val := range tt.values {
				err := s.Set(val)
				if tt.wantErr {
					assert.Error(t, err)
				}
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, s.values)
		})
	}
}
