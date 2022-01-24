package main

import (
	"context"
	"fmt"
	"github.com/NFT-com/indexer/dispatch/local"
	"github.com/NFT-com/indexer/function/ethereum/mainnet/cryptokitties"
	"github.com/NFT-com/indexer/function/ethereum/mainnet/erc1155"
	"github.com/NFT-com/indexer/function/ethereum/mainnet/erc721"
	"github.com/NFT-com/indexer/store/mock"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/function"
	"github.com/NFT-com/indexer/function/ethereum/mainnet"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("failure: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	// Signal catching for clean shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Command line parameter initialization.
	var (
		flagAddress  string
		flagLogLevel string
	)

	pflag.StringVarP(&flagAddress, "address", "a", ":8081", "listening address")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")

	pflag.Parse()

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		return err
	}
	log = log.Level(level)

	mockStore := mock.New()
	dispatcher := local.New("http://localhost:8081/%s")

	var (
		mainnetHandler       = mainnet.New(mockStore, dispatcher)
		erc721Handler        = erc721.New(mockStore)
		erc1155Handler       = erc1155.New(mockStore)
		cryptokittiesHandler = cryptokitties.New(mockStore)
	)

	functionMapper := map[string]function.Function{
		mainnet.Name:       mainnetHandler.Handle,
		erc721.Name:        erc721Handler.Handle,
		erc1155.Name:       erc1155Handler.Handle,
		cryptokitties.Name: cryptokittiesHandler.Handle,
	}

	e := echo.New()
	e.POST("/:function", func(c echo.Context) error {
		functionName := c.Param("function")

		evt := event.Event{}
		if err = c.Bind(&evt); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		functionHandler, ok := functionMapper[functionName]
		if !ok {
			return c.String(http.StatusNotFound, functionName)
		}

		err := functionHandler(context.Background(), &evt)
		if err != nil {
			log.Error().Err(err).Str("function", functionName).Msg("failed to parse event")
		}

		return c.String(http.StatusAccepted, functionName)
	})

	failed := make(chan error)
	done := make(chan struct{})
	go func() {
		log.Info().Msg("launching server")
		if err := e.Start(flagAddress); err != nil {
			log.Error().Err(err).Msg("failed to start server")
		}
		log.Info().Msg("stopping server")
		close(done)
	}()

	select {
	case <-sig:
		if err := e.Close(); err != nil {
			failed <- err
		}
	case err := <-failed:
		return err
	}

	go func() {
		<-sig
		log.Fatal().Msg("forced interruption")
	}()

	<-done

	return nil
}
