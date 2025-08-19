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
	EthGasLimit          uint64 = 21000
	TokenGasLimit        uint64 = 120000
	maxFeePerGas                = big.NewInt(2900000000)
	maxPriorityFeePerGas        = big.NewInt(2600000000)
)

type Withdraw struct {
	db             *database.DB
	chainConf      *config.ChainConfig
	client         node.SolanaClient
	resourceCtx    context.Context
	resourceCancel context.CancelFunc
	tasks          tasks.Group
}

func NewWithdraw(cfg *config.Config, db *database.DB, client node.SolanaClient, shutdown context.CancelCauseFunc) (*Withdraw, error) {
	resCtx, resCancel := context.WithCancel(context.Background())
	return &Withdraw{
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

func (w *Withdraw) Close() error {
	var result error
	w.resourceCancel()
	if err := w.tasks.Wait(); err != nil {
		result = errors.Join(result, fmt.Errorf("failed to await deposit %w"), err)
	}
	return nil
}

func (w *Withdraw) Start() error {
	log.Info("start withdraw......")
	tickerWithdrawsWorker := time.NewTicker(time.Second * 5)
	w.tasks.Go(func() error {
		for range tickerWithdrawsWorker.C {
			log.Info("start solana withdraw......")
		}
		return nil
	})
	return nil
}
