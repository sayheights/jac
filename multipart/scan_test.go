package multipart

import (
	refl "reflect"
	"testing"
)

type scannerNoSkipDef struct {
	ID    string `multi:"id"`
	Count int    `multi:"count"`
	File  string `multi:"file,text/csv"`
}

var scannerNoSkip = scannerNoSkipDef{
	ID:    "ID",
	Count: 10,
	File:  "file.csv",
}

type scannerTagSyntaxDef struct {
	ID string `multi:"id, textttt/csvv"`
}

var scannerTagSyntax = scannerTagSyntaxDef{
	ID: "ID",
}

type scannerSkipAllDef struct {
	NoSkipID string
	scannerNoSkipDef
}

type scannerSkipLastDef struct {
	ID     string `multi:"id"`
	Count  int
	Amount int
}

var scannerSkipLast = &scannerSkipLastDef{
	ID:     "abc",
	Count:  10,
	Amount: 10,
}

var scannerSkipAll = &scannerSkipAllDef{
	NoSkipID:         "NoSkipID",
	scannerNoSkipDef: scannerNoSkip,
}

func Test_scanner_init(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name:    "scanner",
			input:   &scannerNoSkip,
			wantErr: false,
		},
		{
			name:    "scanner",
			input:   scannerNoSkip,
			wantErr: false,
		},
		{
			name:    "scanner",
			input:   "string_input",
			wantErr: true,
		},
		{
			name:    "scanner",
			input:   "int_input",
			wantErr: true,
		},
		{
			name:    "scanner",
			input:   5,
			wantErr: true,
		},
		{
			name:    "scanner",
			input:   []string{"a", "b"},
			wantErr: true,
		},
		{
			name:    "scanner",
			input:   nil,
			wantErr: true,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &scanner{}
			if err := s.init(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("%d: scanner.init() error = %v, wantErr %v",
					i, err, tt.wantErr)
			}
		})
	}
}

func Test_validate(t *testing.T) {
	tests := []struct {
		name    string
		i       interface{}
		wantErr bool
	}{
		{
			name:    "scanner",
			i:       scannerNoSkip,
			wantErr: false,
		},
		{
			name:    "scanner",
			i:       &scannerNoSkip,
			wantErr: false,
		},
		{
			name:    "scanner",
			i:       nil,
			wantErr: true,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validate(tt.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("%d: validate() error = %v, wantErr %v",
					i, err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.Kind() != refl.Struct {
				t.Errorf("%d: validate() value = %v, want reflect.Struct",
					i, got.Kind())
			}
		})
	}
}

func Test_scanner_nextPart(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		index   int
		want    *formPart
		wantErr bool
	}{
		{
			input: scannerSkipAll,
			index: 0,
			want: &formPart{
				name: "id",
				typ:  refl.TypeOf(""),
				val:  refl.ValueOf("ID"),
			},
		},
		{
			input: scannerSkipAll,
			index: 1,
			want: &formPart{
				name: "count",
				typ:  refl.TypeOf(1),
				val:  refl.ValueOf(10),
			},
		},
		{
			input: scannerSkipAll,
			index: 3,
			want: &formPart{
				name: "file",
				typ:  refl.TypeOf(""),
				val:  refl.ValueOf("file.csv"),
				mime: TextCSV,
			},
		},
		{
			input: scannerSkipAll,
			index: 10,
			want:  endPart,
		},
		{
			input: scannerSkipLast,
			index: 2,
			want:  endPart,
		},
		{
			input:   scannerTagSyntax,
			index:   0,
			want:    nil,
			wantErr: true,
		},
		{
			input:   &scannerTagSyntax,
			index:   0,
			want:    nil,
			wantErr: true,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &scanner{}
			err := s.init(tt.input)
			if err != nil {
				t.Errorf("%d: scanner.nextPart() error = %v, wantErr %v",
					i, err, tt.wantErr)
				return
			}
			s.index = tt.index
			_, _ = s.nextPart()
			got, err := s.nextPart()
			if (err != nil) != tt.wantErr {
				t.Errorf("%d: scanner.nextPart() error = %v, wantErr %v",
					i, err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("%d: scanner.nextPart() = %#v, want %#v",
					i, got, tt.want)
			}
		})
	}
}
