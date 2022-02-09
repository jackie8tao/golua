package token

import "testing"

func TestPos_IsValid(t *testing.T) {
	type fields struct {
		Filename string
		Offset   int
		Line     int
		Column   int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"invalid pos", fields{}, false},
		{"valid pos", fields{"", 0, 1, 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pos{
				Filename: tt.fields.Filename,
				Offset:   tt.fields.Offset,
				Line:     tt.fields.Line,
				Column:   tt.fields.Column,
			}
			if got := p.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPos_String(t *testing.T) {
	type fields struct {
		Filename string
		Offset   int
		Line     int
		Column   int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"nil", fields{}, "-"},
		{"file+line+column", fields{"a", 0, 1, 1}, "a:1:1"},
		{"file+line", fields{"a", 0, 1, 0}, "a:1"},
		{"line+column", fields{"", 0, 1, 1}, "1:1"},
		{"line", fields{"", 0, 1, 0}, "1"},
		{"file", fields{"a", 0, 0, 0}, "a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pos{
				Filename: tt.fields.Filename,
				Offset:   tt.fields.Offset,
				Line:     tt.fields.Line,
				Column:   tt.fields.Column,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
