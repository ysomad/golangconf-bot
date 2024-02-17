package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v3"

	mock "github.com/ysomad/golangconf-bot/internal/bot/mock"
)

func Test_middlewareAdminOnly(t *testing.T) {
	type args struct {
		admins []int64
	}
	tests := map[string]struct {
		args      args
		ctxChatID int64 // chat id in context
		wantErr   bool
	}{
		"multiple_admins_ok": {
			args:      args{admins: []int64{1337, 35, 7774, 3461, 2}},
			ctxChatID: 2,
			wantErr:   false,
		},
		"one_admin_ok": {
			args:      args{admins: []int64{1337}},
			ctxChatID: 1337,
			wantErr:   false,
		},
		"zero_chat_id_error": {
			args:      args{admins: []int64{1337}},
			ctxChatID: 0,
			wantErr:   true,
		},
		"negative_chat_id_error": {
			args:      args{admins: []int64{1337}},
			ctxChatID: -87654,
			wantErr:   true,
		},
		"user_not_admin_error": {
			args:      args{admins: []int64{98765}},
			ctxChatID: 1337,
			wantErr:   true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := mock.NewTelebotContext(t)
			c.EXPECT().Chat().Return(&tele.Chat{ID: tt.ctxChatID})

			mw := middlewareAdminOnly(tt.args.admins)

			err := mw(func(tele.Context) error { return nil })(c)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
