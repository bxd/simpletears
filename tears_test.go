package simpletears

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var count int = 0
var l_c sync.Mutex

func Test_Tear(t *testing.T) {
	fmt.Println("testing")
	tt := Default(5).Delay(200, 2000)
	tt.Shed()
	for i := 0; i < 100; i++ {
		j := randJob()
		//		log.Println("")
		tt.AddJobSync(j)
	}
	tt.WipeOff()
	tt.Drying()
}

func result(id string, e error, b []byte) {
	l_c.Lock()
	defer l_c.Unlock()
	count++
	fmt.Println(count, ":\t", e, "\t", len(b))
}

func randJob() *Job {
	j := new(Job).UrlString("http://www.baidu.com/s?wd=" + randString())
	j.SetMethod("GET").SetCallback(result)
	return j
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, rand.Intn(10)+1)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
