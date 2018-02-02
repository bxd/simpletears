package simpletears

import (
	"time"
)

type CalDelay func() time.Duration

type Worm struct {
	cal_delay CalDelay
	quit      chan bool
	busy      bool
	retired   bool
}

func (me *Worm) DelayFunc(f CalDelay) *Worm {
	me.cal_delay = f
	return me
}

func (me *Worm) IsWorking() bool {
	//return (me.quit != nil)
	return me.busy
}

func (me *Worm) DoneFlag() chan bool {
	//	log.Println("in DoneFlag, quit is :", me.quit)
	return me.quit
}

func (me *Worm) Wriggle(j chan *Job) {
	if j == nil {
		return
	}
	go func() {
		me.busy = true
		for {
			if me.retired {
				return
			}
			select {
			case job := <-j:
				job.Exec()
				// case <-me.quit:
				// 	close(me.quit)
				// 	me.quit = nil
				// 	return
			default:
				if !me.busy {
					me.quit <- true
					me.retired = true
				}
			}
			if me.cal_delay != nil {
				if delay := me.cal_delay(); delay > 0 {
					time.Sleep(delay)
				}
			}

		}
	}()
}

func (me *Worm) Retire(immediate bool) {
	me.busy = false
	if immediate {
		me.retired = true
	} else {
		me.quit = make(chan bool)
	}
}
