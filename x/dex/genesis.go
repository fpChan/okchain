package dex

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ordertypes "github.com/okex/okchain/x/order/types"
)

// GenesisState - all slashing state that must be provided at genesis
type GenesisState struct {
	Params        Params                    `json:"params"`
	TokenPairs    []*TokenPair              `json:"token_pairs"`
	WithdrawInfos WithdrawInfos             `json:"withdraw_infos"`
	ProductLocks  ordertypes.ProductLockMap `json:"product_locks"`
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
// TODO: check how the added params' influence export facility
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:        *DefaultParams(),
		TokenPairs:    nil,
		WithdrawInfos: nil,
		ProductLocks:  *ordertypes.NewProductLockMap(),
	}
}

// ValidateGenesis validates the slashing genesis parameters
func ValidateGenesis(data GenesisState) error {
	return nil
}

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper IKeeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)

	// reset token pair
	for _, pair := range data.TokenPairs {
		err := keeper.SaveTokenPair(ctx, pair)
		if err != nil {
			panic(err)
		}
	}

	// reset delay withdraw queue
	for _, withdrawInfo := range data.WithdrawInfos {
		keeper.SetWithdrawInfo(ctx, withdrawInfo)
		keeper.SetWithdrawCompleteTimeAddress(ctx, withdrawInfo.CompleteTime, withdrawInfo.Owner)
	}

	for k, v := range data.ProductLocks.Data {
		keeper.LockTokenPair(ctx, k, v)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper IKeeper) (data GenesisState) {
	params := keeper.GetParams(ctx)
	tokenPairs := keeper.GetTokenPairsFromStore(ctx)
	var withdrawInfos WithdrawInfos
	keeper.IterateWithdrawInfo(ctx, func(_ int64, withdrawInfo WithdrawInfo) (stop bool) {
		withdrawInfos = append(withdrawInfos, withdrawInfo)
		return false
	})
	return GenesisState{
		Params:        params,
		TokenPairs:    tokenPairs,
		WithdrawInfos: withdrawInfos,
		ProductLocks:  *keeper.LoadProductLocks(ctx),
	}
}
