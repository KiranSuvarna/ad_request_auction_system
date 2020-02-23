package misc

import "time"

// Ad auction request
type AuctionReq struct {
	AuctionID string `json:"auction_id"`
}

// Bidder registration request
type BidderRequest struct {
	Name  string        `json:"name"`
	Delay time.Duration `json:"delay"`
}

// Bidder details
type Bidder struct {
	ID    string        `json:"id"`
	Name  string        `json:"name"`
	Host  string        `json:"host"`
	Delay time.Duration `json:"delay"`
}

// Bid response details
type BidResponse struct {
	BidderID string  `json:"bidder_id"`
	Amount   float64 `json:"amount"`
}

// Meta holds the request status
type Meta struct {
	Status  int    `json:"status_code"`
	Message string `json:"error_message,omitempty"`
}

type Response struct {
	Data interface{} `json:"data,omitempty"`
	Meta Meta        `json:"meta"`
}
