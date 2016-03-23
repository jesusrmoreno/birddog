package gordita

import "github.com/jesusrmoreno/birddog/Godeps/_workspace/src/github.com/siddontang/ledisdb/ledis"

// Errors
const (
	qKey       = "___o_o___"
	qSortedKey = "sorted____o_o___"
	qExistsKey = "___=_=___"
	ErrDequeue = "There was an error when popping from Queue"
	ErrEnqueue = "There was an error when pushing to Queue"
)

// Q is a simple Q
type Q struct {
	DB *ledis.DB
}

// PQ ...
type PQ struct {
	DB *ledis.DB
}

// PQMember ...
type PQMember struct {
	Value    []byte
	Priority int64
}

// New returns and initializes a new queue
func New(db *ledis.DB) *Q {
	return &Q{
		DB: db,
	}
}

// NewPQ returns a new priority Queue
func NewPQ(db *ledis.DB) *PQ {
	return &PQ{
		DB: db,
	}
}

// IsEmpty returns whether the priority queue is empty
func (q *PQ) IsEmpty() bool {
	return q.Size() == 0
}

// Empty empties the priority queue
func (q *PQ) Empty() {
	q.DB.ZClear([]byte(qSortedKey))
}

// Push pushed to the priority queue
func (q *PQ) Push(val []byte, score int64) {
	sp := ledis.ScorePair{
		Member: val,
		Score:  score,
	}
	q.DB.ZAdd([]byte(qSortedKey), sp)
}

// Pop gets the highest priority item
func (q *PQ) Pop() PQMember {
	scorePair, _ := q.DB.ZRange([]byte(qSortedKey), 0, 0)
	item := PQMember{}
	if len(scorePair) > 0 {
		item.Priority = scorePair[0].Score
		item.Value = scorePair[0].Member
	}
	return item
}

// Size returns the size of the priority Queue
func (q *PQ) Size() int64 {
	size, _ := q.DB.ZCard([]byte(qSortedKey))
	return size
}

// Size returns the size of the queue ...
func (q *Q) Size() int64 {
	s, _ := q.DB.LLen([]byte(qKey))
	return s
}

// IsEmpty returns whether or not the q is empty
func (q *Q) IsEmpty() bool {
	return q.Size() == 0
}

// Empty empties the queue ...
func (q *Q) Empty() int64 {
	elem, _ := q.DB.LClear([]byte(qKey))
	q.DB.SClear([]byte(qExistsKey))
	return elem
}

// Push pushes a new value into the queue
func (q *Q) Push(val []byte) {
	q.DB.LPush([]byte(qKey), val)
	q.DB.SAdd([]byte(qExistsKey), val)
}

// Pop removes a value from the Queue
func (q *Q) Pop() []byte {
	val, _ := q.DB.RPop([]byte(qKey))
	q.DB.SRem([]byte(qExistsKey), val)
	return val
}

// Contains returns whether a value is already in the queue.
func (q *Q) Contains(val []byte) bool {
	exists, _ := q.DB.SIsMember([]byte(qExistsKey), val)
	return exists != 0
}
