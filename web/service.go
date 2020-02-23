package web

import (
	"sync"
	"time"

	"bitbucket.org/greedygames/ad_request_auction_system/misc"
	st "bitbucket.org/greedygames/ad_request_auction_system/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Service HTTP server info
type Service struct {
	refreshInterval time.Duration
	shutdownChan    chan bool
	domain          string

	router *gin.Engine
	wg     sync.WaitGroup

	AppName   string
	Version   string
	BuildTime string
}

var StoreCon *st.Conn
var AuctioneerStore *st.Conn
var BidderStore *st.Conn

// NewService Create a new service
func NewService(conf *misc.Config) (*Service, error) {

	Init()

	s := &Service{
		refreshInterval: conf.RefreshInterval,
		router:          gin.New(),
		shutdownChan:    make(chan bool),
		domain:          conf.HTTP.Domain,
	}

	s.router.Use(gin.Logger())

	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s.router.GET("/", s.index)
	s.router.GET("/ping", s.ping)

	v1 := s.router.Group("/v1")
	{
		v1.POST("/auction", s.auctionHandler)
		v1.POST("/bidder/register", s.registerBidder)
		v1.GET("/bidder/all", s.getBidders)
	}

	return s, nil
}

// Start the web service
func (s *Service) Start(address string) error {
	return s.router.Run(address)
}

// Close all threads and free up resources
func (s *Service) Close() {
	close(s.shutdownChan)

	s.wg.Wait()

}

func Init() {
	StoreCon = st.NewStore()
	AuctioneerStore = StoreCon
	BidderStore = StoreCon
}
