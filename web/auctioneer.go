package web

import (
	"bitbucket.org/greedygames/ad_request_auction_system/misc"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sort"
	"sync"
	"time"
)

type bidds []*misc.BidResponse

func (a bidds) Len() int           { return len(a) }
func (a bidds) Less(i, j int) bool { return a[i].Amount > a[j].Amount }
func (a bidds) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (s *Service) auctionHandler(c *gin.Context) {
	var (
		input misc.AuctionReq
		wg    sync.WaitGroup
	)

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bidders := AuctioneerStore.Bidder().List()
	if len(bidders) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no bidders available"})
		return
	}
	wg.Add(len(bidders))

	data := make(chan *misc.BidResponse, len(bidders))
	for _, b := range bidders {
		go collectBidResponse(input.AuctionID, b.Host, &wg, data)
	}

	var bidRes bidds
	for i := 0; i < len(bidders); i++ {
		if d := <-data; d != nil {
			bidRes = append(bidRes, d)
		}
	}

	wg.Wait()
	close(data)
	sort.Sort(bidRes)
	if len(bidRes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bidders not responding within time"})
		return
	}

	s.responseWriter(c, bidRes[0], http.StatusOK)
}

func collectBidResponse(auctionID, host string, wg *sync.WaitGroup,
	data chan *misc.BidResponse) {
	var err error
	body := bytes.NewBuffer(nil)
	json.NewEncoder(body).Encode(map[string]interface{}{
		"auction_id": auctionID,
	})

	defer func() {
		wg.Done()
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
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var res struct {
		Data *misc.BidResponse `json:"data"`
		Meta misc.Meta         `json:"meta"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return
	}

	data <- res.Data
}
