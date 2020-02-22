package web

import (
	"sync"
	"time"

	"bitbucket.org/youplus/dashboard_server/misc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

)

// Service HTTP server info
type Service struct {
	refreshInterval time.Duration
	shutdownChan    chan bool

	router *gin.Engine
	wg     sync.WaitGroup

	AppName   string
	Version   string
	BuildTime string
}

// NewService Create a new service
func NewService(conf *misc.Config) (*Service, error) {

	s := &Service{
		refreshInterval: conf.RefreshInterval,
		router:          gin.New(),
		shutdownChan:    make(chan bool),
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
	v1.Use(s.authMiddleware())
	{

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
