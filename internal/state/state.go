package state

import (
	"fmt"
	"log/slog"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

type Step uint8

const (
	StepUnknown Step = iota
	StepUploadingSchedule
)

func key(tguid int64) string {
	return fmt.Sprintf("state:%d", tguid)
}

type State struct {
	Data []byte
	Step Step
}

type storage struct {
	lru *expirable.LRU[string, State]
}

func NewStorage(lru *expirable.LRU[string, State]) storage {
	return storage{lru: lru}
}

func (s storage) Save(tguid int64, st State) {
	_ = s.lru.Add(key(tguid), st) // omit checking is key evicted, it must be handled in EvictedCallback
}

func (s storage) Get(tguid int64) State {
	if st, ok := s.lru.Get(key(tguid)); ok {
		return st
	}

	return State{}
}

func (s storage) Del(tguid int64) {
	if present := s.lru.Remove(key(tguid)); !present {
		slog.Warn("trying to remove not existing state", "tg_uid", tguid)
	}
}
