package announcer_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/tvanriper/announcer"
)

func Example() {
	sender := announcer.New(5)
	listenA := sender.Listen()
	listenB := sender.Listen()
	sender.Send("Hullo")
	defer sender.Close()
	if la, ok := <-listenA.Listen(); ok {
		fmt.Println(la.(string))
	}
	if lb, ok := <-listenB.Listen(); ok {
		fmt.Println(lb.(string))
	}

	// Output:
	// Hullo
	// Hullo
}

func IsSame(l interface{}, r interface{}) bool {
	if ls, ok := l.(string); ok {
		if rs, ok := r.(string); ok {
			return ls == rs
		} else {
			return false
		}
	}
	if li, ok := l.(int); ok {
		if ri, ok := r.(int); ok {
			return li == ri
		} else {
			return false
		}
	}

	if la, ok := l.([]string); ok {
		if ra, ok := r.([]string); ok {
			if len(la) == len(ra) {
				for i := 0; i < len(la); i++ {
					if la[i] != ra[i] {
						return false
					}
				}
				return true
			}
		}
		return false
	}

	if lf, ok := l.(float64); ok {
		if rf, ok := r.(float64); ok {
			return lf == rf
		}
	}
	return false
}

func IsStop(c interface{}) bool {
	if ch, ok := c.(string); ok {
		if ch == "stop" {
			return true
		}
	}
	return false
}

func TestAnnouncements(t *testing.T) {
	// Strategy:
	// Create three listeners.
	// Two listeners will record what they hear into an array.
	// The third listener, a 'stray', will test what they hear in the same
	// thread as the sender.
	// We will send different data types, to include an array.
	// One of the listeners recording into an array will stop listening when it
	// hears the string "stop".
	// The other listener recording into an array will record everything up to
	// the announcer's close.
	// The stray will also stop listening before the last item.
	// We will also attempt to send something after the Announcer is closed,
	// and we will create a new listener (and listen on it) after the Announcer
	// is closed.
	storagea := make([]interface{}, 0)
	storageb := make([]interface{}, 0)

	starting := sync.WaitGroup{}
	starting.Add(2)

	stopping := sync.WaitGroup{}
	stopping.Add(2)

	a := announcer.New(2)

	List := func(l []interface{}, s string) {
		for i := 0; i < len(l); i++ {
			t.Logf("s[%d]: %v", i, l[i])
		}
	}

	Compare := func(l []interface{}, r []interface{}, s string) bool {
		if len(l) != len(r) {
			t.Errorf("%s: expected %d items but got %d", s, len(l), len(r))
		} else {
			for i := 0; i < len(l); i++ {
				if !IsSame(l[i], r[i]) {
					List(r, s)
					t.Errorf("%s: expected %v but got %v", s, l[i], r[i])
					return false
				}
			}
		}
		return true
	}

	go func() {
		l := a.Listen()
		starting.Done()
		for {
			data, ok := <-l.Listen()
			if ok {
				storagea = append(storagea, data)
				if IsStop(data) {
					break
				}
			} else {
				break
			}
		}
		stopping.Done()
	}()

	go func() {
		l := a.Listen()
		starting.Done()
		for {
			data, ok := <-l.Listen()
			if ok {
				storageb = append(storageb, data)
			} else {
				break
			}
		}
		stopping.Done()
	}()

	stray := a.Listen()

	// Ensure the listeners have been created and are listening.
	starting.Wait()

	expects := make([]interface{}, 0)
	expects = append(expects, "Hello")
	expects = append(expects, 1)
	expects = append(expects, []string{"this", "can", "be", "interesting"})
	expects = append(expects, "stop")

	// Send the data.
	for _, expect := range expects {
		err := a.Send(expect)
		if err != nil {
			t.Errorf("expected no error but received %s", err)
		}
		got := <-stray.Listen()
		if !IsSame(expect, got) {
			t.Errorf("expected %v but got %v", expect, got)
		}
	}
	// Closing listener before announcer is closed.
	stray.Close()

	// Testing that stray doesn't get the announcement
	go func() {
		if got, ok := <-stray.Listen(); ok {
			t.Errorf("should not have received %v", got)
		}
	}()

	lastOne := 3.141592653598
	// Waiting, just a bit, to ensure the stray is listening if it's going to.
	time.Sleep(10 * time.Millisecond)

	a.Send(lastOne)

	// storagea should not have received the last item, as it was told to quit
	// if it saw "stop"
	if !Compare(expects, storagea, "storagea") {
		List(storagea, "storagea")
	}
	a.Close()

	// Ensure the listeners have finished, and received everything.
	stopping.Wait()

	// storageb should have received the last item.
	expects = append(expects, lastOne)

	if !Compare(expects, storageb, "storageb") {
		List(storageb, "storageb")
	}

	// Test that we generate an error when sending after a close
	err := a.Send(a)
	if err == nil {
		t.Errorf("expected an error but send returned nil")
	}

	// Test that a listener cannot hear anything after closing the announcer.
	wha := a.Listen()
	if _, ok := <-wha.Listen(); ok {
		t.Errorf("lister's channel should have been closed")
	}
}
