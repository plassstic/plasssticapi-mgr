package bot

import "plassstic.tech/gopkg/plassstic-mgr/lib/ent"

func InitPresenceThreads(c *ent.Client) {
	repo = &goroutinesRepo{threads: make(map[int]chan struct{})}
	repo.populate(c)
}

func HasThread(userId int64) bool {
	_, ok := repo.threads[int(userId)]
	return ok
}

func StopPresenceThreads() {
	for _, closeChan := range repo.threads {
		closeChan <- struct{}{}
	}
}
