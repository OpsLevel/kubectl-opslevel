package signal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

var (
	onlyOneSignalHandler = make(chan struct{})
	stopChannel          = make(chan struct{})
)

// InitSignalHandler
// Usage:
//
//		func Start() {
//	   log.Info().Msg("Starting...")
//	   <-opslevel_common.InitSignalHandler() // Block until signals
//	   log.Info().Msg("Stopping...")
//		}
func InitSignalHandler() <-chan struct{} {
	close(onlyOneSignalHandler) // panics when called twice

	closeChannel := make(chan os.Signal, 2)
	signal.Notify(closeChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-closeChannel
		close(stopChannel)
		<-closeChannel
		os.Exit(1) // second signal. Exit directly.
	}()

	return stopChannel
}

func Run(name string) {
	log.Info().Msgf("Starting %s ...", name)
	<-InitSignalHandler() // Block until signals
	log.Info().Msgf("Stopping %s ...", name)
}

func Stop() {
	close(stopChannel)
}
