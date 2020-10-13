package wechat_test

import (
	"testing"

	"github.com/sohaha/wechat"
)

func TestWechat(t *testing.T) {
	wx := wechat.New(&wechat.Mp{
		AppID:     "wx9d1fcb71007a71b0",
		AppSecret: "c4132441ded3301bda2d2373609959e1",
	})
	t.Log(wx.GetAccessToken())
}
