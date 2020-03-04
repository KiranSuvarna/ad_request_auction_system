package web

import (
	"fmt"
	"net/http"

	"bitbucket.org/greedygames/ad_request_auction_system/misc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

func (s *Service) getBidders(c *gin.Context) {
	var(
		bidders []*misc.Bidder
		ba      []string
		err error
	)

	if ba, err = s.rc.Keys(fmt.Sprintf("gg_*")); err == nil {
		if biddersMeta, err := s.rc.GetMultiple(ba); err == nil {	
			for k, v := range biddersMeta {
				var bidder *misc.Bidder
				if err := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(v, &bidder); err != nil {
					log.WithError(err).WithField("key", k).Warn("Failed to decode cache value")
				}
				bidders = append(bidders, bidder)
			}
			s.responseWriter(c, bidders, http.StatusCreated)
		}	
	} else {
		log.WithError(err).WithField("key", fmt.Sprintf("gg_bidder")).Warn("Failed to get bidders count from cache")	
	}
	
	s.responseWriter(c, "No bidders available", http.StatusNoContent)
}

func (s *Service) registerBidder(c *gin.Context) {
	var req misc.BidderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	bidder := &misc.Bidder{
		ID:    fmt.Sprintf("bidder_%s", uuid.New().String()),
		Name:  req.Name,
		Host:  fmt.Sprintf("http://%s", c.Request.Host),
		Delay: req.Delay,
	}

	if err := s.rc.Set(fmt.Sprintf("gg_%s", bidder.ID), bidder); err != nil {
		log.WithError(err).WithField("key", fmt.Sprintf("gg_%s", bidder.ID)).Warn("Failed to set cache")
	}

	res := &misc.Response{
		Data: bidder,
		Meta: misc.Meta{
			Status: http.StatusCreated,
		},
	}

	s.responseWriter(c, res, http.StatusCreated)
}
