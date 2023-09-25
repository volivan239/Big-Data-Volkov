package main

import "time"

func snapshotManager() {
	for ; ; time.Sleep(time.Minute) {
		mutex.Lock()

		previousSnapshot = snapshot
		previousJournal = journal

		if len(journal) > 0 {
			snapshot = journal[len(journal)-1]
			journal = nil
		}

		mutex.Unlock()
	}
}
