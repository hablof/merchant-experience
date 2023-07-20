package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_contains(t *testing.T) {

	tests := []struct {
		name  string
		slice []uint64
		elem  uint64
		want  bool
	}{
		{
			name:  "первый элемент",
			slice: []uint64{0, 1, 2, 3, 4, 5},
			elem:  0,
			want:  true,
		},
		{
			name:  "последний элемент",
			slice: []uint64{0, 1, 2, 3, 4, 5, 6},
			elem:  6,
			want:  true,
		},
		{
			name:  "центральный элемент",
			slice: []uint64{0, 1, 2, 3, 4, 5, 6},
			elem:  3,
			want:  true,
		},
		{
			name:  "отсутствующий элемент в середине",
			slice: []uint64{0, 1, 2, 4, 5, 6},
			elem:  3,
			want:  false,
		},
		{
			name:  "отсутствующий элемент больше большего",
			slice: []uint64{0, 1, 2, 3, 4, 5, 6},
			elem:  7,
			want:  false,
		},
		{
			name:  "отсутствующий элемент меньше меньшего",
			slice: []uint64{1, 2, 3, 4, 5, 6},
			elem:  0,
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := contains(tt.slice, tt.elem)
			assert.Equal(t, tt.want, got)
		})
	}
}
