package store

import (
	"bitbucket.org/greedygames/ad_request_auction_system/misc"
)

// Store global store interface
type Store interface {
	Bidder() Bidders
}

// Bidders store interface
type Bidders interface {
	Add(bidder *misc.Bidder)
	List() []*misc.Bidder
	Count() int
}

// Conn struct holds the store connection
type Conn struct {
	BidderConn Bidders
}

// NewStore inits new store connection
func NewStore() *Conn {
	// new allocates zeroed storage for a new item or type whatever and then returns a pointer to it
	// same as  conn := Conn{}
	conn := new(Conn)
	conn.BidderConn = NewBidderStore(conn)

	return conn
}

// Bidder implements the store interface
func (s *Conn) Bidder() Bidders {
	return s.BidderConn
}
