//
// Copyright 2020 Wireline, Inc.
//

package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/wirelineio/wns/x/auction/internal/types"
)

// prefixIDToAuctionIndex is the prefix for ID -> Auction index in the KVStore.
// Note: This is the primary index in the system.
// Note: Golang doesn't support const arrays.
var prefixIDToAuctionIndex = []byte{0x00}

// prefixOwnerToAuctionsIndex is the prefix for the Owner -> [Auction] index in the KVStore.
var prefixOwnerToAuctionsIndex = []byte{0x01}

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper
	supplyKeeper  supply.Keeper

	// Track auction usage in other cosmos-sdk modules (more like a usage tracker).
	usageKeepers []types.AuctionUsageKeeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.

	paramstore params.Subspace
}

// AuctionClientKeeper is the subset of functionality exposed to other modules.
type AuctionClientKeeper interface {
	HasAuction(ctx sdk.Context, id types.ID) bool
	GetAuction(ctx sdk.Context, id types.ID) types.Auction
	MatchAuctions(ctx sdk.Context, matchFn func(*types.Auction) bool) []*types.Auction
}

// NewKeeper creates new instances of the auction Keeper
func NewKeeper(accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, supplyKeeper supply.Keeper,
	usageKeepers []types.AuctionUsageKeeper, storeKey sdk.StoreKey, cdc *codec.Codec, paramstore params.Subspace) Keeper {
	return Keeper{
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		supplyKeeper:  supplyKeeper,
		usageKeepers:  usageKeepers,
		storeKey:      storeKey,
		cdc:           cdc,
		paramstore:    paramstore.WithKeyTable(ParamKeyTable()),
	}
}

// Generates Auction ID -> Auction index key.
func getAuctionIndexKey(id types.ID) []byte {
	return append(prefixIDToAuctionIndex, []byte(id)...)
}

// Generates Owner -> Auctions index key.
func getOwnerToAuctionsIndexKey(owner sdk.AccAddress, auctionID types.ID) []byte {
	return append(append(prefixOwnerToAuctionsIndex, []byte(owner)...), []byte(auctionID)...)
}

// SaveAuction - saves a auction to the store.
func (k Keeper) SaveAuction(ctx sdk.Context, auction types.Auction) {
	store := ctx.KVStore(k.storeKey)

	// Auction ID -> Auction index.
	store.Set(getAuctionIndexKey(auction.ID), k.cdc.MustMarshalBinaryBare(auction))

	// Owner -> [Auction] index.
	store.Set(getOwnerToAuctionsIndexKey(auction.OwnerAddress, auction.ID), []byte{})
}

// HasAuction - checks if a auction by the given ID exists.
func (k Keeper) HasAuction(ctx sdk.Context, id types.ID) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(getAuctionIndexKey(id))
}

// DeleteAuction - deletes the auction.
func (k Keeper) DeleteAuction(ctx sdk.Context, auction types.Auction) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(getAuctionIndexKey(auction.ID))
	store.Delete(getOwnerToAuctionsIndexKey(auction.OwnerAddress, auction.ID))
}

// GetAuction - gets a record from the store.
func (k Keeper) GetAuction(ctx sdk.Context, id types.ID) types.Auction {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(getAuctionIndexKey(id))
	var obj types.Auction
	k.cdc.MustUnmarshalBinaryBare(bz, &obj)

	return obj
}

// ListAuctions - get all auctions.
func (k Keeper) ListAuctions(ctx sdk.Context) []types.Auction {
	var auctions []types.Auction

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, prefixIDToAuctionIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj types.Auction
			k.cdc.MustUnmarshalBinaryBare(bz, &obj)
			auctions = append(auctions, obj)
		}
	}

	return auctions
}

// QueryAuctionsByOwner - query auctions by owner.
func (k Keeper) QueryAuctionsByOwner(ctx sdk.Context, ownerAddress string) []types.Auction {
	var auctions []types.Auction

	ownerPrefix := append(prefixOwnerToAuctionsIndex, []byte(ownerAddress)...)
	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, ownerPrefix)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		auctionID := itr.Key()[len(ownerPrefix):]
		bz := store.Get(append(prefixIDToAuctionIndex, auctionID...))
		if bz != nil {
			var obj types.Auction
			k.cdc.MustUnmarshalBinaryBare(bz, &obj)
			auctions = append(auctions, obj)
		}
	}

	return auctions
}

// MatchAuctions - get all matching auctions.
func (k Keeper) MatchAuctions(ctx sdk.Context, matchFn func(*types.Auction) bool) []*types.Auction {
	var auctions []*types.Auction

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, prefixIDToAuctionIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj types.Auction
			k.cdc.MustUnmarshalBinaryBare(bz, &obj)
			if matchFn(&obj) {
				auctions = append(auctions, &obj)
			}
		}
	}

	return auctions
}

// CreateAuction creates a new auction.
func (k Keeper) CreateAuction(ctx sdk.Context, msg types.MsgCreateAuction) (*types.Auction, sdk.Error) {
	// Might be called from another module directly, always validate.
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	// Generate auction ID.
	account := k.accountKeeper.GetAccount(ctx, msg.Signer)
	if account == nil {
		return nil, sdk.ErrInvalidAddress("Account not found.")
	}

	auctionID := types.AuctionID{
		Address:  msg.Signer,
		AccNum:   account.GetAccountNumber(),
		Sequence: account.GetSequence(),
	}.Generate()

	// Take fees from account.
	totalFee := msg.CommitFee.Add(msg.RevealFee)
	sdkErr := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Signer, types.ModuleName, sdk.NewCoins(totalFee))
	if err != nil {
		return nil, sdkErr
	}

	// Compute timestamps.
	now := ctx.BlockTime()
	commitsEndTime := now.Add(time.Duration(msg.CommitsDuration) * time.Second)
	revealsEndTime := now.Add(time.Duration(msg.CommitsDuration+msg.RevealsDuration) * time.Second)

	auction := types.Auction{
		ID:             types.ID(auctionID),
		Status:         types.AuctionStatusCommitPhase,
		OwnerAddress:   msg.Signer,
		CreateTime:     now,
		CommitsEndTime: commitsEndTime,
		RevealsEndTime: revealsEndTime,
		CommitFee:      msg.CommitFee,
		RevealFee:      msg.RevealFee,
		MinimumBid:     msg.MinimumBid,
	}

	// Save auction in store.
	k.SaveAuction(ctx, auction)

	return &auction, nil
}

// GetAuctionModuleBalances gets the auction module account(s) balances.
func (k Keeper) GetAuctionModuleBalances(ctx sdk.Context) map[string]sdk.Coins {
	balances := map[string]sdk.Coins{}
	accountNames := []string{types.ModuleName}

	for _, accountName := range accountNames {
		moduleAddress := k.supplyKeeper.GetModuleAddress(accountName)
		moduleAccount := k.accountKeeper.GetAccount(ctx, moduleAddress)
		if moduleAccount != nil {
			balances[accountName] = moduleAccount.GetCoins()
		}
	}

	return balances
}
