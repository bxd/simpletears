package simpletears

import (
	"log"
	"math/rand"
	"runtime"
	"time"
)

type Tear struct {
	min_delay int
	max_delay int
	worm_num  uint
	jobs      Jobs
	worms     []*Worm
}

//create a new crawler with default config
func Default(num uint) *Tear {
	if num == 0 {
		return nil
	}
	t := &Tear{
		// min_delay: 1 * time.Second,
		// max_delay: 1 * time.Second,
		worm_num: num,
	}
	t.jobs.job_buf = make(chan *Job, t.worm_num*10)
	t.worms = make([]*Worm, t.worm_num)
	for k := range t.worms {
		t.worms[k] = new(Worm)
	}
	runtime.SetFinalizer(t, t.Makeup)
	return t
}

func (me *Tear) Makeup(o *Tear) {
	close(me.jobs.job_buf)
	me.jobs.job_buf = nil
}

func (me *Tear) WipeOff() {
	for _, v := range me.worms {
		//log.Printf("Before Retire:\t %+v\n", v)
		v.Retire(false)
		//log.Printf("After Retire:\t %+v\n", v)
	}
}

func (me *Tear) Hire(num uint) *Tear {
	me.worm_num = num
	for _, v := range me.worms {
		v.Retire(true)
	}
	me.worms = make([]*Worm, me.worm_num)
	for k := range me.worms {
		me.worms[k] = new(Worm)
	}
	return me
}

func (me *Tear) Shed() {
	for _, v := range me.worms {
		if v.IsWorking() {
			continue
		}
		v.DelayFunc(me.CalDelay).Wriggle(me.jobs.job_buf)
	}
}

func (me *Tear) Drying() {
	for _, v := range me.worms {
		q := v.DoneFlag()
		if q != nil {
			<-q
			close(q)
		} else {
			log.Println("Quit flag  is nil!!!", q)
		}
	}
}

//set max and min delay between twice wriggles
func (me *Tear) Delay(min, max int) *Tear {
	if min >= max {
		return me
	}
	me.min_delay = min
	me.max_delay = max
	return me
}

//set a steady delay between twice wriggles
func (me *Tear) SteadyDelay(delay int) *Tear {
	me.min_delay = delay
	me.max_delay = me.min_delay
	return me
}

func (me *Tear) CalDelay() time.Duration {
	if me.max_delay == 0 {
		return 0
	}
	if me.min_delay == me.max_delay {
		return time.Duration(me.max_delay)
	}

	rand.Seed(time.Now().UnixNano())
	ra := rand.Intn(me.max_delay - me.min_delay)
	return time.Duration(me.min_delay + ra)
}

func (me *Tear) AddJobSync(j *Job) *Tear {
	if j == nil {
		return me
	}
	me.jobs.Push(j)
	return me
}

func (me *Tear) AddJobAsync(j *Job) *Tear {
	if j == nil {
		return me
	}
	go me.jobs.Push(j)
	return me
}

// func (me *Tear) AddUrlJob(u *url.URL) *Tear {
// 	j := new(Job).Url(*u)
// 	return me.AddJobSync(j)
// }

// func (me *Tear) AddUrlStringJob(str string) *Tear {
// 	j := new(Job).UrlString(str)
// 	return me.AddJobSync(j)
// }
