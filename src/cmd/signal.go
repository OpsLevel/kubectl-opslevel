package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2024"
	"github.com/rs/zerolog/log"
)

func InitSignalHandler(parent context.Context, queue chan<- opslevel_jq_parser.ServiceRegistration) context.Context {
	ctx, cancel := context.WithCancel(parent)
	closeChannel := make(chan os.Signal, 1)
	signal.Notify(closeChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-closeChannel
		log.Info().Str("signal", sig.String()).Msg("Handling interruption")
		cancel()
		close(queue)
	}()
	return ctx
}
