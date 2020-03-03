package main

import (
	"bytes"
	"encoding/json"
	errs "errors"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"bitbucket.org/greedygames/ad_request_auction_system/misc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	// Port in which the bidder client to run
	Port int
	// Delay after the bidder to respond to auction request
	Delay time.Duration
	// BidderID unique identifier of bidder
	BidderID string
)

func main() {
	name := flag.String("name", "Sample", "Bidder's name")
	port := flag.Int("port", 0, "Bidder's port")
	delay := flag.Uint("delay", 0, "Bidder's response delay")

	flag.Parse()

	if *port == 0 {
		log.Fatal("invalid port or port required")
	}

	if *delay > 500 {
		log.Println("delay more than 500ms")
	}

	Port = *port
	Delay = (time.Duration(*delay) * time.Millisecond)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("v1/bid", bidHandler)

	// register with auctioneer
	if err := registerWithAuctioneer(*name); err != nil {
		log.Fatal("Can't able to register with auctioneer. Err: ", err.Error())
	}

	fmt.Printf("Starting bidder %s on port :%d with delay %d\n",
		*name, *port, *delay)

	router.Run(fmt.Sprintf(":%d", *port))
}

func bidHandler(c *gin.Context) {
	<-time.After(Delay)
	rand.Seed(time.Now().UnixNano())
	min := 100
	max := 3000

	data := map[string]interface{}{
		"bidder_id": BidderID,
		"amount":    rand.Intn(max-min+1) + min,
	}

	c.JSON(http.StatusOK, &misc.Response{
		Data: data,
		Meta: misc.Meta{Status: http.StatusOK},
	})
}

func registerWithAuctioneer(name string) error {
	url := fmt.Sprintf("%s%s", "http://localhost:5000", "/v1/bidder/register")
	body := bytes.NewBuffer(nil)
	json.NewEncoder(body).Encode(map[string]interface{}{
		"name":  name,
		"delay": Delay,
	})

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}
	req.Host = fmt.Sprintf("%s:%d", "localhost", Port)
	req.Header.Set("Content-Type", "application/json")

	if resp, err := http.DefaultClient.Do(req); err == nil {
		// Closed the body to avoid the leaking
		defer resp.Body.Close()

		// This checks for non 200 responses also
		if resp.StatusCode != 201 {
			return errs.New("Server error!!")
		}

		var res struct {
			Data *misc.Bidder `json:"data"`
			Meta misc.Meta    `json:"meta"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return err
		}

		if resp.StatusCode > 201 {
			return errs.New(res.Meta.Message)
		}

		BidderID = res.Data.ID
	} else {
		// Error is not nil when client's CheckRedirect func failed or if there there are any HTTP protocol error
		return err
	}
	return nil
}
