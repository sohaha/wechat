package wechat_test

import (
	"testing"

	"github.com/sohaha/wechat"
)

func TestWechat(t *testing.T) {
	wx := wechat.New(&wechat.Mp{
		AppID:     "wx6a24b584b45b6791",
		AppSecret: "cc31573bfa7af4cdc2ba327357af9234",
	})
	t.Log(wx.GetAccessToken())
}
