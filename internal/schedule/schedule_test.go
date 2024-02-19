package schedule

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_readCSV(t *testing.T) {
	type args struct {
		f io.Reader
	}
	tests := map[string]struct {
		args    args
		want    [][]string
		wantErr bool
	}{
		"nil_reader": {
			args:    args{f: nil},
			want:    nil,
			wantErr: true,
		},
		"invalid_csv_extra_col": {
			args: args{f: strings.NewReader(`Start (MSK Time Zone),Duration (min),Title,Speakers,URL
2024-02-19 10:00,60,Introduction to Go,Gopher,"https://example.com/intro-to-go",Extra Field`)},
			want:    nil,
			wantErr: true,
		},
		"ok": {
			args: args{f: strings.NewReader(`Start (MSK Time Zone),Duration (min),Title,Speakers,URL
2024-02-19 10:00,60,Introduction to Go,Gopher,"https://example.com/intro-to-go"`)},
			want: [][]string{
				{"Start (MSK Time Zone)", "Duration (min)", "Title", "Speakers", "URL"},
				{"2024-02-19 10:00", "60", "Introduction to Go", "Gopher", "https://example.com/intro-to-go"},
			},
			wantErr: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := readCSV(tt.args.f)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestNewFromCSV(t *testing.T) {
	type args struct {
		file io.Reader
	}
	tests := map[string]struct {
		args         args
		wantSchedule []Talk
		wantErrs     []error
	}{
		"error_nil_reader": {
			args:         args{file: nil},
			wantSchedule: nil,
			wantErrs:     []error{errNilReader},
		},
		"error_invalid_csv": {
			args: args{file: strings.NewReader(`Start (MSK Time Zone),Duration (min),Title,Speakers,URL
2024-02-19 10:00,60,Introduction to Go,Gopher,"https://example.com/intro-to-go",Extra Field`)},
			wantSchedule: nil,
			wantErrs:     []error{errInvalidCSV},
		},
		"error_empty_cols": {
			args: args{file: strings.NewReader(`Start (MSK Time Zone),Duration (min),Title,Speakers,URL
21/07/2024 10:00:00,60,Introduction to Go,,https://example.com/intro-to-go
21/07/2024 10:00:00,60,,test,
,60,Introduction to Go,test,https://example.com/intro-to-go
21/07/2024 10:00:00,60,Introduction to Go,test,https://example.com/intro-to-go
21/07/2024 10:00:00,60,,test,
21/07/2024 10:00:00,,Introduction to Go,test,https://example.com/intro-to-go
,60,Introduction to Go,test,https://example.com/intro-to-go
21/07/2024 10:00:00,60,Introduction to Go,test,https://example.com/intro-to-go
21/07/2024 10:00:00,60,Introduction to Go,test,https://example.com/intro-to-go`)},
			wantSchedule: nil,
			wantErrs: []error{
				CSVError{
					cause: errEmptyCol,
					row:   2,
					col:   4,
				},
				CSVError{
					cause: errEmptyCol,
					row:   3,
					col:   3,
				},
				CSVError{
					cause: errEmptyCol,
					row:   3,
					col:   5,
				},
				CSVError{
					cause: errEmptyCol,
					row:   4,
					col:   1,
				},
				CSVError{
					cause: errEmptyCol,
					row:   6,
					col:   3,
				},
				CSVError{
					cause: errEmptyCol,
					row:   6,
					col:   5,
				},
				CSVError{
					cause: errEmptyCol,
					row:   7,
					col:   2,
				},
				CSVError{
					cause: errEmptyCol,
					row:   8,
					col:   1,
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, errs := NewFromCSV(tt.args.file)

			assert.Equal(t, len(tt.wantErrs), len(errs))

			if len(tt.wantErrs) > 0 {
				for i, wantErr := range tt.wantErrs {
					assert.ErrorIs(t, errs[i], wantErr)

					// test if error occured in specific row and col
					if csvErr, ok := errs[i].(CSVError); ok {
						assert.Equal(t, wantErr, csvErr)
					}
				}
			}

			assert.Equal(t, tt.wantSchedule, schedule)
		})
	}
}
