package slogx

import (
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseLevel(t *testing.T) {
	type args struct {
		level string
	}
	tests := map[string]struct {
		args args
		want slog.Level
	}{
		"empty": {
			args: args{level: ""},
			want: slog.LevelInfo,
		},
		"invalid": {
			args: args{level: "INVALID_LOG_LEVEL"},
			want: slog.LevelInfo,
		},
		"info1": {
			args: args{level: "info"},
			want: slog.LevelInfo,
		},
		"info2": {
			args: args{level: "InFo"},
			want: slog.LevelInfo,
		},
		"warn1": {
			args: args{level: "warn"},
			want: slog.LevelWarn,
		},
		"warn2": {
			args: args{level: "wARN"},
			want: slog.LevelWarn,
		},
		"warn3": {
			args: args{level: "warning"},
			want: slog.LevelWarn,
		},
		"warn4": {
			args: args{level: "WarNiNG"},
			want: slog.LevelWarn,
		},
		"error1": {
			args: args{level: "error"},
			want: slog.LevelError,
		},
		"error2": {
			args: args{level: "ErRoR"},
			want: slog.LevelError,
		},
		"error3": {
			args: args{level: "err"},
			want: slog.LevelError,
		},
		"error4": {
			args: args{level: "ERr"},
			want: slog.LevelError,
		},
		"debug1": {
			args: args{level: "debug"},
			want: slog.LevelDebug,
		},
		"debug2": {
			args: args{level: "DebUG"},
			want: slog.LevelDebug,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, parseLevel(tt.args.level))
		})
	}
}

func TestNewHandler(t *testing.T) {
	type args struct {
		w         io.Writer
		format    string
		level     string
		addSource bool
	}
	tests := map[string]struct {
		args args
		want slog.Handler
	}{
		"json_handler": {
			args: args{
				w:         os.Stdout,
				format:    "json",
				level:     "warning",
				addSource: true,
			},
			want: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelWarn,
			}),
		},
		"text_handler": {
			args: args{
				w:         os.Stderr,
				format:    "text",
				level:     "error",
				addSource: false,
			},
			want: slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				AddSource: false,
				Level:     slog.LevelError,
			}),
		},
		"empty_format": {
			args: args{
				w:         os.Stderr,
				format:    "",
				level:     "error",
				addSource: false,
			},
			want: slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				AddSource: false,
				Level:     slog.LevelError,
			}),
		},
		"unsupported_format": {
			args: args{
				w:         os.Stdout,
				format:    "xml",
				level:     "info",
				addSource: true,
			},
			want: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelInfo,
			}),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewHandler(tt.args.w, tt.args.format, tt.args.level, tt.args.addSource))
		})
	}
}
