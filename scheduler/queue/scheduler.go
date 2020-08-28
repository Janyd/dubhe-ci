package queue

import "dubhe-ci/core"

type scheduler struct {
	*queue
	*canceller
}

func New(store core.BuildStore) core.Scheduler {
	return &scheduler{
		queue:     newQueue(store),
		canceller: newCanceller(),
	}
}
