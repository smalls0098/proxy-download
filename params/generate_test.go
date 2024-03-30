package params

import (
	"testing"
	"time"
)

func TestGen(t *testing.T) {
	gen, err := Gen("http://127.0.0.1:13822/proxy", "1.mp4", Params{
		Tag:    1,
		Url:    "https://scontent.cdninstagram.com/v/t51.29350-15/408431926_3450217495240215_7454044517912965548_n.jpg?stp=dst-jpg_e35_p1080x1080&_nc_ht=scontent.cdninstagram.com&_nc_cat=103&_nc_ohc=M8xadq-7rIgAX-rOgCV&edm=APs17CUBAAAA&ccb=7-5&oh=00_AfCFcNGV7FPbKlX2Pe1FM93hcbmpYVHv5t4fUixlWQNsdQ&oe=660D887E&_nc_sid=10d13b",
		Expire: time.Now().Unix() + 7200,
	}, "smalls0098")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gen)
}
