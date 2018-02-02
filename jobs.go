package simpletears

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type JobCallback func(string, error, []byte)
type JobSetting struct {
	http_method string
	header      map[string]string
	call_back   JobCallback
}

//describe a job
type Job struct {
	id      string
	u       string
	body    io.Reader
	setting JobSetting
}

func (me *Job) Url(uu url.URL) *Job {
	me.u = uu.String()
	return me
}

func (me *Job) UrlString(u_str string) *Job {
	me.u = u_str
	return me
}

func (me *Job) SetMethod(m string) *Job {
	switch strings.ToUpper(m) {
	case "GET":
		me.setting.http_method = "GET"
	case "POST":
		me.setting.http_method = "POST"
	}
	return me
}

func (me *Job) SetHeader(h map[string]string) *Job {
	me.setting.header = h
	return me
}

func (me *Job) SetCallback(f JobCallback) *Job {
	me.setting.call_back = f
	return me
}

func (me *Job) Exec() {
	req, e := http.NewRequest(me.setting.http_method,
		me.u,
		me.body,
	)
	if e != nil {
		log.Println(e)
		me.setting.call_back(me.id, e, nil)
		return
	}
	for k, v := range me.setting.header {
		req.Header.Set(k, v)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, e := client.Do(req)
	if e != nil {
		log.Println(e)
		me.setting.call_back(me.id, e, nil)
		return
	}
	defer resp.Body.Close()
	if b, e := ioutil.ReadAll(resp.Body); e == nil {
		if me.setting.call_back != nil {
			me.setting.call_back(me.id, nil, b)
		}
	} else {
		log.Println(e)
		me.setting.call_back(me.id, e, nil)
		return
	}
}

/*
type Jobs struct {
	head *Job
	tail *Job
	l    sync.Mutex
}

//push a job into the queue
func (me *Jobs) Push(j *Job) {
	if j == nil {
		return
	}
	me.l.Lock()
	defer me.l.Unlock()
	if me.tail == nil {
		me.head = j
		me.tail = j
		return
	}
	me.tail.next = j
	me.tail = j
}

//pop a job from the queue
func (me *Jobs) Pop() *Job {
	if me.head == nil {
		return nil
	}
	me.l.Lock()
	defer me.l.Unlock()
	ret := me.head
	me.head = me.head.next
	if me.head == nil {
		me.tail = nil
	}
	ret.next = nil
	return ret
}

//if the jobs queue is empty
func (me *Jobs) IsEmpty() bool {
	me.l.Lock()
	defer me.l.Unlock()
	if me.head == nil {
		return true
	}
	return false
}
*/

type Jobs struct {
	job_buf chan *Job
}

func (me *Jobs) Push(j *Job) {
	if j == nil {
		return
	}

	if me.job_buf != nil {
		me.job_buf <- j
	} else {
		log.Println("job buf is nil")
	}
}
