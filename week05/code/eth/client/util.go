// Copied from perun-examples/payment-channel/client/util.go (Apache-2.0, PolyCrypt GmbH).


package client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethchannel "github.com/perun-network/perun-eth-backend/channel"
	ethwallet "github.com/perun-network/perun-eth-backend/wallet"
	swallet "github.com/perun-network/perun-eth-backend/wallet/simple"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/wire"
)

const (
	txFinalityDepth = 1 // Number of blocks required to confirm a transaction.
)

// CreateContractBackend creates a new contract backend.
func CreateContractBackend(
	nodeURL string,
	chainID uint64,
	w *swallet.Wallet,
) (ethchannel.ContractBackend, error) {
	signer := types.LatestSignerForChainID(new(big.Int).SetUint64(chainID))
	transactor := swallet.NewTransactor(w, signer)

	ethClient, err := ethclient.Dial(nodeURL)
	if err != nil {
		return ethchannel.ContractBackend{}, err
	}

	return ethchannel.NewContractBackend(ethClient, ethchannel.MakeChainID(big.NewInt(int64(chainID))), transactor, txFinalityDepth), nil
}

// WalletAddress returns the wallet address of the client.
func (c *PaymentClient) WalletAddress() common.Address {
	return common.Address(*c.account[1].(*ethwallet.Address))
}

// WireAddress returns the wire address of the client.
func (c *PaymentClient) WireAddress() map[wallet.BackendID]wire.Address {
	return c.waddress
}

// EthToWei converts a given amount in ETH to Wei.
func EthToWei(ethAmount *big.Float) (weiAmount *big.Int) {
	weiPerEth := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	weiPerEthFloat := new(big.Float).SetInt(weiPerEth)
	weiAmountFloat := new(big.Float).Mul(ethAmount, weiPerEthFloat)
	weiAmount, _ = weiAmountFloat.Int(nil)
	return weiAmount
}

// WeiToEth converts a given amount in Wei to ETH.
func WeiToEth(weiAmount *big.Int) (ethAmount *big.Float) {
	weiPerEth := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	weiPerEthFloat := new(big.Float).SetInt(weiPerEth)
	weiAmountFloat := new(big.Float).SetInt(weiAmount)
	return new(big.Float).Quo(weiAmountFloat, weiPerEthFloat)
}
