package web

import (
	"fmt"
	"net/http"

	"bitbucket.org/greedygames/ad_request_auction_system/misc"
	"github.com/gin-gonic/gin"
)

func (s *Service) getBidders(c *gin.Context) {
	s.responseWriter(c, BidderStore.Bidder().List(), http.StatusOK)
}

func (s *Service) registerBidder(c *gin.Context) {
	var req misc.BidderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	bidder := &misc.Bidder{
		ID:    fmt.Sprintf("bidder_%d", BidderStore.Bidder().Count()+1),
		Name:  req.Name,
		Host:  fmt.Sprintf("http://%s", c.Request.Host),
		Delay: req.Delay,
	}

	res := &misc.Response{
		Data: bidder,
		Meta: misc.Meta{
			Status: http.StatusCreated,
		},
	}

	BidderStore.Bidder().Add(bidder)

	s.responseWriter(c, res, http.StatusCreated)
}
