package treepath

// Fifo is a basic FIFO queue based on a circular list that resizes as needed.
// https://play.golang.org/p/m15vAaFQ9r
type Fifo struct {
	nodes []interface{}
	head  int
	tail  int
	count int
}

func NewFifo(size int) *Fifo {
	if size <= 0 {
		size = 1
	}
	q := &Fifo{nodes: make([]interface{}, size)}
	return q
}

// Push adds a node to the queue.
func (q *Fifo) Push(n interface{}) {
	if q.head == q.tail && q.count > 0 {
		nodes := make([]interface{}, len(q.nodes)*2)
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.head])
		q.head = 0
		q.tail = len(q.nodes)
		q.nodes = nodes
	}
	q.nodes[q.tail] = n
	q.tail = (q.tail + 1) % len(q.nodes)
	q.count++
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Fifo) Pop() interface{} {
	if q.count == 0 {
		return nil
	}
	node := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return node
}

func (q *Fifo) Len() int {
	return q.count
}
