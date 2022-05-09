package pkg

import "sync"

var locks sync.Map

func IsFree(name, namespace string) bool {
	_, ok := locks.Load(name + "|" + namespace)

	return !ok
}

func Lock(name, namespace string) {
	locks.Store(name+"|"+namespace, "LOCKED")
}

func UnLock(name, namespace string) {
	locks.Delete(name + "|" + namespace)
}
