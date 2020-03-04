package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"bitbucket.org/greedygames/ad_request_auction_system/misc"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

type bidds []*misc.BidResponse

func (a bidds) Len() int           { return len(a) }
func (a bidds) Less(i, j int) bool { return a[i].Amount > a[j].Amount }
func (a bidds) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (s *Service) auctionHandler(c *gin.Context) {
	var (
		input misc.AuctionReq
		ba []string
		bidders []*misc.Bidder
		err error
	)

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ba, err = s.rc.Keys(fmt.Sprintf("gg_*")); err != nil {
		log.WithError(err).WithField("key", fmt.Sprintf("gg_bidder")).Warn("Failed to get bidders count from cache")
	} 

	if len(ba) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no bidders available"})
		return
	}

	if biddersMeta, err := s.rc.GetMultiple(ba); err == nil {
		for k, v := range biddersMeta {
			var bidder *misc.Bidder
			if err := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(v, &bidder); err != nil {
				log.WithError(err).WithField("key", k).Warn("Failed to decode cache value")
			}
			bidders = append(bidders, bidder)
		}
	}

	data := make(chan *misc.BidResponse, len(bidders))
	for _, b := range bidders {
		go collectBidResponse(input.AuctionID, b.Host, data)
	}

	var bidRes bidds
	for i := 0; i < len(bidders); i++ {
		if d := <-data; d != nil {
			bidRes = append(bidRes, d)
		}
	}

	close(data)
	sort.Sort(bidRes)
	if len(bidRes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bidders not responding within time"})
		return
	}

	s.responseWriter(c, bidRes[0], http.StatusOK)
}

func collectBidResponse(auctionID, host string,
	data chan *misc.BidResponse) {
	var err error
	body := bytes.NewBuffer(nil)
	json.NewEncoder(body).Encode(map[string]interface{}{
		"auction_id": auctionID,
	})

	defer func() {
		if err != nil {
			log.WithError(err).Warn("no data")
			data <- nil
		}
	}()

	url := fmt.Sprintf("%s%s", host, "/v1/bid")
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 190*time.Millisecond)
	defer cancel()

	if resp, err := http.DefaultClient.Do(req.WithContext(ctx)); err == nil {
		// Closing the body to avoid the leaking
		defer resp.Body.Close()

		var res struct {
			Data *misc.BidResponse `json:"data"`
			Meta misc.Meta         `json:"meta"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return
		}

		data <- res.Data
	} else {
		data <- nil
		log.WithError(err).Warn("Failed to get the response from bidder")
		return
	}	
}
