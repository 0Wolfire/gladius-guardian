package guardian

import (
	"container/list"
	"sync"
)

// FixedSizeLog is a log storage that only keeps a max number of entries, and
// deletes old ones
type FixedSizeLog struct {
	logList    *list.List // Linked list for efficient popping of old elements
	maxLogSize int        // How many lines can our log be before we delete old lines
	mux        sync.Mutex
}

// NewFixedSizeLog returns a new FixedSizeLog with the specified max size of
// log entries to keep
func NewFixedSizeLog(maxSize int) *FixedSizeLog {
	return &FixedSizeLog{
		logList:    list.New(),
		maxLogSize: maxSize,
		mux:        sync.Mutex{},
	}
}

// Append adds to the log
func (fsl *FixedSizeLog) Append(line string) {
	fsl.mux.Lock()
	defer fsl.mux.Unlock()

	// Delete the first element if our list grows too long
	if fsl.logList.Len() >= fsl.maxLogSize {
		fsl.logList.Remove(fsl.logList.Front())
	}
	fsl.logList.PushBack(line) // Always add the line to the log
}

// LogLines returns a string slice representing the underlying values
func (fsl *FixedSizeLog) LogLines() []string {
	toReturn := make([]string, 0, 1000)
	for e := fsl.logList.Front(); e != nil; e = e.Next() {
		toReturn = append(toReturn, e.Value.(string))
	}
	return toReturn
}
