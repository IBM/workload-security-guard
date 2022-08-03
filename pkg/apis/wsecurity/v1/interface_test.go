package v1

import "testing"

func TestPile_Clear(t *testing.T) {

	tests := []struct {
		name string
	}{
		{"Clear"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pile{}
			p.Clear()
		})
	}
}
