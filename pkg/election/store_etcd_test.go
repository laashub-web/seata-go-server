package election

import (
	"context"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func TestAddExpectLeader(t *testing.T) {
	cli, err := clientv3.NewFromURL("http://127.0.0.1:2379")
	if err != nil {
		t.Errorf("clientv3 failed with %+v", err)
		return
	}

	opts := options{}
	opts.client = cli
	opts.leaseSec = 1

	s := &store{
		opts:          opts,
		client:        opts.client,
		lessor:        clientv3.NewLease(opts.client),
		leasors:       make(map[uint64]clientv3.LeaseID),
		watcheCancels: make(map[uint64]context.CancelFunc),
		watchers:      make(map[uint64]clientv3.Watcher),
	}

	c := make(chan int, 1)
	leaderFunc := func() {
		c <- 1
	}

	go func() {
		err := s.campaignLeader(1, 10, leaderFunc, func() {})
		if err != nil {
			t.Errorf("campaignLeader failed with %+v", err)
			return
		}
	}()

	select {
	case <-c:
		err := s.addExpectLeader(1, 10, 100)
		if err != nil {
			t.Errorf("campaignLeader failed with %+v", err)
			return
		}

		break
	case <-time.After(time.Second * 10):
		t.Errorf("campaignLeader timeout")
		return
	}
}
