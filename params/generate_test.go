package params

import (
	"testing"
	"time"
)

func TestGen(t *testing.T) {
	gen, err := Gen("http://127.0.0.1:13822/proxy", "1.mp4", Params{
		Tag:    1,
		Url:    "https://v3-default.365yg.com/af147f1c864883759ad26401086de850/65b60371/video/tos/cn/tos-cn-v-0015c001/oECEAAwUhfAgn1AgEemzC8Kr8kmAxEyvoAgtBP/?a=1398&br=3064&bt=3064&btag=e00028000&cd=0%7C0%7C0%7C0&ch=0&cr=0&cv=1&dr=0&ds=3&dy_q=1706423423&dy_va_biz_cert=&er=6&ft=3JiezGTo62DXhNvjVi1N5dzfu17LYQIkaYlc&l=20240128143023EBE9308CBFC676A7B22A&lr=watermark&mime_type=video_mp4&net=5&qs=13&rc=MzdkOzc6ZnhzcDMzNGkzM0BpMzdkOzc6ZnhzcDMzNGkzM0Beay5kcjRfZWtgLS1kLS9zYSNeay5kcjRfZWtgLS1kLS9zcw%3D%3D&logo_type=douyin",
		Expire: time.Now().Unix() + 7200,
	}, "smalls0098")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gen)
}
