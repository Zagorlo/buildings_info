package serviceCache

import (
	"buildings_info/consts"
	"buildings_info/logging"
	"buildings_info/models"
	"context"
	"encoding/json"
	"sync"
)

type JsonOrderedCache struct {
	items [][]byte
	l     sync.RWMutex
}

func NewJsonOrderedCache(cfg models.BuildingsCache) *JsonOrderedCache {
	return &JsonOrderedCache{
		items: make([][]byte, 0, cfg.Size),
	}
}

func (oc *JsonOrderedCache) lock() {
	oc.l.Lock()
}

func (oc *JsonOrderedCache) unlock() {
	oc.l.Unlock()
}

func (oc *JsonOrderedCache) PrependItems(requestCtx context.Context, items []models.CacheItem) {
	if len(items) == 0 {
		return
	}

	oc.lock()
	defer oc.unlock()

	var pool [][]byte
	for i := range items {
		bytes, err := json.Marshal(items[i])
		if err != nil {
			logging.NewErrorContainer(requestCtx, err, consts.JsonMarshalErrorStatus)

			continue
		}

		pool = append(pool, bytes)
	}

	pool = append(pool, oc.items...)

	if len(pool) < 10 {
		oc.items = pool
	} else if len(pool) > 0 {
		oc.items = pool[:10]
	}
}

func (oc *JsonOrderedCache) RetrieveItems() []byte {
	oc.lock()
	defer oc.unlock()

	answer := []byte("[")

	for i := range oc.items {
		answer = append(answer, oc.items[i]...)
		if i < len(oc.items)-1 {
			answer = append(answer, []byte(",")...)
		}
	}

	answer = append(answer, []byte("]")...)

	return answer
}
