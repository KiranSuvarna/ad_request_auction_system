package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bitbucket.org/greedygames/ad_request_auction_system/misc"
	"bitbucket.org/greedygames/ad_request_auction_system/web"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// AppVersion is the application version
	AppVersion = "1.0"
	// Version is the git commit version (set by Makefile)
	Version = "none"
	// BuildTime application build time (set by Makefile)
	BuildTime = "none"

	version = flag.Bool("version", false, "print version string")

	appName = "gg_ad_auction_system"
)

func main() {
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	fullVersion := fmt.Sprintf("%s-%s", AppVersion, Version)

	if *version {
		fmt.Printf("%s v%s (%s)\n", appName, fullVersion, BuildTime)
		flag.PrintDefaults()

		return
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %v", err)
	}

	var config misc.Config
	err := viper.Unmarshal(&config)

	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	misc.InitLogging(&config.Log)

	log.WithFields(log.Fields{
		"app":       appName,
		"version":   Version,
		"buildTime": BuildTime,
	}).Info("Starting up")

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	s, err := web.NewService(&config)
	if err != nil {
		log.WithError(err).Error("Failed to start service")
		return
	}

	s.AppName = appName
	s.Version = Version
	s.BuildTime = BuildTime

	go func() {
		if err := s.Start(config.HTTP.Address); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start web service")
		}

		termChan <- syscall.SIGTERM
	}()

	select {
	case <-termChan:
		if s != nil {
			s.Close()
		}

		os.Exit(0)
	}
}
