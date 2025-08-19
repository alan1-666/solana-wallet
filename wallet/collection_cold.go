package wallet

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"
	
	"github.com/ethereum/go-ethereum/log"

	"github.com/alan1-666/solana-wallet/common/tasks"
	"github.com/alan1-666/solana-wallet/config"
	"github.com/alan1-666/solana-wallet/database"
	"github.com/alan1-666/solana-wallet/wallet/node"
)

var (
	CollectionFunding = big.NewInt(10000000000000000)
	ColdFunding       = big.NewInt(2000000000000000000)
)

type CollectionCold struct {
	db             *database.DB
	chainConf      *config.ChainConfig
	client         node.SolanaClient
	resourceCtx    context.Context
	resourceCancel context.CancelFunc
	tasks          tasks.Group
}

func NewCollectionCold(cfg *config.Config, db *database.DB, client node.SolanaClient, shutdown context.CancelCauseFunc) (*CollectionCold, error) {
	resCtx, resCancel := context.WithCancel(context.Background())
	return &CollectionCold{
		db:             db,
		chainConf:      &cfg.Chain,
		client:         client,
		resourceCtx:    resCtx,
		resourceCancel: resCancel,
		tasks: tasks.Group{HandleCrit: func(err error) {
			shutdown(fmt.Errorf("critical error in deposit: %w", err))
		}},
	}, nil
}

func (cc *CollectionCold) Close() error {
	var result error
	cc.resourceCancel()
	if err := cc.tasks.Wait(); err != nil {
		result = errors.Join(result, fmt.Errorf("failed to await deposit %w"), err)
	}
	return nil
}

func (cc *CollectionCold) Start() error {
	log.Info("start collection and cold......")
	tickerCollectionColdWorker := time.NewTicker(time.Second * 5)
	cc.tasks.Go(func() error {
		for range tickerCollectionColdWorker.C {
			log.Info("solana collection......")

		}
		return nil
	})

	cc.tasks.Go(func() error {
		for range tickerCollectionColdWorker.C {
			log.Info("solana cold......")

		}
		return nil
	})

	return nil
}
