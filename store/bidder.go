package store

import "bitbucket.org/greedygames/ad_request_auction_system/misc"

var bidders = make([]*misc.Bidder, 0)

// BidderStore implements the Store interface
type BidderStore struct {
	*Conn
}

// NewBidderStore returns new store object
func NewBidderStore(st *Conn) *BidderStore {
	return &BidderStore{st}
}

// Add register the new bidder into the list
func (b *BidderStore) Add(bidder *misc.Bidder) {
	bidders = append(bidders, bidder)
}

// List returns all the regisytered bidders
func (b *BidderStore) List() []*misc.Bidder {
	return bidders
}

// Count returns the count of regisytered bidders
func (b *BidderStore) Count() int {
	return len(bidders)
}
