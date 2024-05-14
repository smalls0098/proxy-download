package params

import (
	"testing"
	"time"
)

func TestGen(t *testing.T) {
	gen, err := Gen("http://127.0.0.1:13822/proxy", "1.mp4", Params{
		Tag:    0,
		Url:    "https://rr4---sn-npoe7nes.googlevideo.com/videoplayback?expire=1715672599&ei=t8FCZsbQN4jjjuMP96O7uAw&ip=13.213.8.126&id=o-ANaaEx1JnQ6JF1pSsCrP1kXucDxe9mwjJq5LnlZkv1TY&itag=22&source=youtube&requiressl=yes&xpc=EgVo2aDSNQ%3D%3D&mh=Q-&mm=31%2C29&mn=sn-npoe7nes%2Csn-npoeeney&ms=au%2Crdu&mv=m&mvi=4&pl=21&gcr=sg&initcwndbps=1037500&vprv=1&svpuc=1&mime=video%2Fmp4&rqh=1&cnr=14&ratebypass=yes&dur=29.698&lmt=1715375982182154&mt=1715650618&fvip=1&c=ANDROID&txp=6308224&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cxpc%2Cgcr%2Cvprv%2Csvpuc%2Cmime%2Crqh%2Ccnr%2Cratebypass%2Cdur%2Clmt&sig=AJfQdSswRgIhANOXLln-otlHKJFX95FkIHfs79E5f0FxstN7nBt8C9fcAiEA4gzB-Hqo3gvaVNhTdv8QmWqwhzKWrBpxo-qWE2sTbNg%3D&lsparams=mh%2Cmm%2Cmn%2Cms%2Cmv%2Cmvi%2Cpl%2Cinitcwndbps&lsig=AHWaYeowRgIhAIWv9JufDEaI7tyLssq1trbBpysCrpYN9f24WTQbqHfxAiEAv0uVNSZwjQdJhaxd0XK18qoWP1fNmY4sCRt0W0lWGf8%3D",
		Expire: time.Now().Unix() + 7200,
	}, "smalls0098")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gen)
}
