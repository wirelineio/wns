//
// Copyright 2020 Wireline, Inc.
//

package keeper

import (
	"fmt"
	"net/url"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/go-amino"
	"github.com/wirelineio/wns/x/auction"
	"github.com/wirelineio/wns/x/nameservice/internal/helpers"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

func GetCIDToNamesIndexKey(id types.ID) []byte {
	return append(PrefixCIDToNamesIndex, []byte(id)...)
}

// Generates name -> NameAuthority index key.
func GetNameAuthorityIndexKey(name string) []byte {
	return append(PrefixNameAuthorityRecordIndex, []byte(name)...)
}

// Generates WRN -> NameRecord index key.
func GetNameRecordIndexKey(wrn string) []byte {
	return append(PrefixWRNToNameRecordIndex, []byte(wrn)...)
}

// HasNameAuthority - checks if a name/authority exists.
func (k Keeper) HasNameAuthority(ctx sdk.Context, name string) bool {
	return HasNameAuthority(ctx.KVStore(k.storeKey), name)
}

// HasNameAuthority - checks if a name authority entry exists.
func HasNameAuthority(store sdk.KVStore, name string) bool {
	return store.Has(GetNameAuthorityIndexKey(name))
}

// SetNameAuthority creates the NameAutority record.
func (k Keeper) SetNameAuthority(ctx sdk.Context, name string, authority types.NameAuthority) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetNameAuthorityIndexKey(name), k.cdc.MustMarshalBinaryBare(authority))
	k.updateBlockChangesetForNameAuthority(ctx, name)
}

// GetNameAuthority - gets a name authority from the store.
func GetNameAuthority(store sdk.KVStore, codec *amino.Codec, name string) *types.NameAuthority {
	authorityKey := GetNameAuthorityIndexKey(name)
	if !store.Has(authorityKey) {
		return nil
	}

	bz := store.Get(authorityKey)
	var obj types.NameAuthority
	codec.MustUnmarshalBinaryBare(bz, &obj)

	return &obj
}

// GetNameAuthority - gets a name authority from the store.
func (k Keeper) GetNameAuthority(ctx sdk.Context, name string) *types.NameAuthority {
	return GetNameAuthority(ctx.KVStore(k.storeKey), k.cdc, name)
}

// AddRecordToNameMapping adds a name to the record ID -> []names index.
func AddRecordToNameMapping(store sdk.KVStore, codec *amino.Codec, id types.ID, wrn string) {
	reverseNameIndexKey := GetCIDToNamesIndexKey(id)

	var names []string
	if store.Has(reverseNameIndexKey) {
		codec.MustUnmarshalBinaryBare(store.Get(reverseNameIndexKey), &names)
	}

	nameSet := sliceToSet(names)
	nameSet.Add(wrn)
	store.Set(reverseNameIndexKey, codec.MustMarshalBinaryBare(setToSlice(nameSet)))
}

// RemoveRecordToNameMapping removes a name from the record ID -> []names index.
func RemoveRecordToNameMapping(store sdk.KVStore, codec *amino.Codec, id types.ID, wrn string) {
	reverseNameIndexKey := GetCIDToNamesIndexKey(id)

	var names []string
	codec.MustUnmarshalBinaryBare(store.Get(reverseNameIndexKey), &names)
	nameSet := sliceToSet(names)
	nameSet.Remove(wrn)

	if nameSet.Cardinality() == 0 {
		// Delete as storing empty slice throws error from baseapp.
		store.Delete(reverseNameIndexKey)
	} else {
		store.Set(reverseNameIndexKey, codec.MustMarshalBinaryBare(setToSlice(nameSet)))
	}
}

// SetNameRecord - sets a name record.
func SetNameRecord(store sdk.KVStore, codec *amino.Codec, wrn string, id types.ID, height int64) {
	nameRecordIndexKey := GetNameRecordIndexKey(wrn)

	var nameRecord types.NameRecord
	if store.Has(nameRecordIndexKey) {
		bz := store.Get(nameRecordIndexKey)
		codec.MustUnmarshalBinaryBare(bz, &nameRecord)
		nameRecord.History = append(nameRecord.History, nameRecord.NameRecordEntry)

		// Update old CID -> []Name index.
		if nameRecord.NameRecordEntry.ID != "" {
			RemoveRecordToNameMapping(store, codec, nameRecord.NameRecordEntry.ID, wrn)
		}
	}

	nameRecord.NameRecordEntry = types.NameRecordEntry{
		ID:     id,
		Height: height,
	}

	store.Set(nameRecordIndexKey, codec.MustMarshalBinaryBare(nameRecord))

	// Update new CID -> []Name index.
	if id != "" {
		AddRecordToNameMapping(store, codec, id, wrn)
	}
}

// SetNameRecord - sets a name record.
func (k Keeper) SetNameRecord(ctx sdk.Context, wrn string, id types.ID) {
	SetNameRecord(ctx.KVStore(k.storeKey), k.cdc, wrn, id, ctx.BlockHeight())

	// Update changeset for name.
	k.updateBlockChangesetForName(ctx, wrn)
}

// HasNameRecord - checks if a name record exists.
func (k Keeper) HasNameRecord(ctx sdk.Context, wrn string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetNameRecordIndexKey(wrn))
}

// GetNameRecord - gets a name record from the store.
func GetNameRecord(store sdk.KVStore, codec *amino.Codec, wrn string) *types.NameRecord {
	nameRecordKey := GetNameRecordIndexKey(wrn)
	if !store.Has(nameRecordKey) {
		return nil
	}

	bz := store.Get(nameRecordKey)
	var obj types.NameRecord
	codec.MustUnmarshalBinaryBare(bz, &obj)

	return &obj
}

// GetNameRecord - gets a name record from the store.
func (k Keeper) GetNameRecord(ctx sdk.Context, wrn string) *types.NameRecord {
	return GetNameRecord(ctx.KVStore(k.storeKey), k.cdc, wrn)
}

// ListNameAuthorityRecords - get all name authority records.
func (k Keeper) ListNameAuthorityRecords(ctx sdk.Context) map[string]types.NameAuthority {
	nameAuthorityRecords := make(map[string]types.NameAuthority)

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, PrefixNameAuthorityRecordIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var record types.NameAuthority
			k.cdc.MustUnmarshalBinaryBare(bz, &record)
			nameAuthorityRecords[string(itr.Key()[len(PrefixNameAuthorityRecordIndex):])] = record
		}
	}

	return nameAuthorityRecords
}

// ListNameRecords - get all name records.
func (k Keeper) ListNameRecords(ctx sdk.Context) map[string]types.NameRecord {
	nameRecords := make(map[string]types.NameRecord)

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, PrefixWRNToNameRecordIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var record types.NameRecord
			k.cdc.MustUnmarshalBinaryBare(bz, &record)
			nameRecords[string(itr.Key()[len(PrefixWRNToNameRecordIndex):])] = record
		}
	}

	return nameRecords
}

// ResolveWRN resolves a WRN to a record.
func (k Keeper) ResolveWRN(ctx sdk.Context, wrn string) *types.Record {
	return ResolveWRN(ctx.KVStore(k.storeKey), k.cdc, wrn)
}

// ResolveWRN resolves a WRN to a record.
func ResolveWRN(store sdk.KVStore, codec *amino.Codec, wrn string) *types.Record {
	nameKey := GetNameRecordIndexKey(wrn)

	if store.Has(nameKey) {
		bz := store.Get(nameKey)
		var obj types.NameRecord
		codec.MustUnmarshalBinaryBare(bz, &obj)

		recordExists := HasRecord(store, obj.ID)
		if !recordExists || obj.ID == "" {
			return nil
		}

		record := GetRecord(store, codec, obj.ID)
		return &record
	}

	return nil
}

// UsesAuction returns true if the auction is used for an name authority.
func (k RecordKeeper) UsesAuction(ctx sdk.Context, auctionID auction.ID) bool {
	// TODO(ashwin): Implement auction ID -> NameAuthority index.
	return false
}

// NotifyAuction is called on auction state change.
func (k RecordKeeper) NotifyAuction(ctx sdk.Context, auctionID auction.ID) {
	// TODO(ashwin): Update authority status based on auction status/winner.
}

// ProcessReserveAuthority reserves a name authority.
func (k Keeper) ProcessReserveAuthority(ctx sdk.Context, msg types.MsgReserveAuthority) (string, sdk.Error) {
	wrn := fmt.Sprintf("wrn://%s", msg.Name)

	parsedWRN, err := url.Parse(wrn)
	if err != nil {
		return "", sdk.ErrInternal("Invalid name.")
	}

	name := parsedWRN.Host
	if fmt.Sprintf("wrn://%s", name) != wrn {
		return "", sdk.ErrInternal("Invalid name.")
	}

	// Check if name already reserved.
	if k.HasNameAuthority(ctx, name) {
		return "", sdk.ErrInternal("Name already exists.")
	}

	if strings.Contains(name, ".") {
		return k.ProcessReserveSubAuthority(ctx, name, msg)
	}

	// Reserve name with signer as owner.
	sdkErr := k.createAuthority(ctx, name, msg.Signer, true)
	if sdkErr != nil {
		return "", sdkErr
	}

	return name, nil
}

func (k Keeper) createAuthority(ctx sdk.Context, name string, owner sdk.AccAddress, isRoot bool) sdk.Error {
	ownerAccount := k.accountKeeper.GetAccount(ctx, owner)
	if ownerAccount == nil {
		return sdk.ErrUnknownAddress("Account not found.")
	}

	pubKey := ownerAccount.GetPubKey()
	if pubKey == nil {
		return sdk.ErrInvalidPubKey("Account public key not set.")
	}

	auctionID := auction.ID("")
	status := types.AuthorityActive

	// Create auction if root authority and name auctions are enabled.
	if isRoot && k.NameAuctionsEnabled(ctx) {
		commitFee, err := sdk.ParseCoin(k.NameAuctionCommitFee(ctx))
		if err != nil {
			return sdk.ErrInvalidCoins("Invalid name auction commit fee.")
		}

		revealFee, err := sdk.ParseCoin(k.NameAuctionRevealFee(ctx))
		if err != nil {
			return sdk.ErrInvalidCoins("Invalid name auction reveal fee.")
		}

		minimumBid, err := sdk.ParseCoin(k.NameAuctionMinimumBid(ctx))
		if err != nil {
			return sdk.ErrInvalidCoins("Invalid name auction minimum bid.")
		}

		// Create an auction.
		msg := auction.NewMsgCreateAuction(auction.Params{
			CommitsDuration: k.NameAuctionCommitsDuration(ctx),
			RevealsDuration: k.NameAuctionRevealsDuration(ctx),
			CommitFee:       commitFee,
			RevealFee:       revealFee,
			MinimumBid:      minimumBid,
		}, owner)

		// TODO(ashwin): Perhaps consume extra gas for auction creation.
		auction, sdkErr := k.auctionKeeper.CreateAuction(ctx, msg)
		if sdkErr != nil {
			return sdkErr
		}

		// TODO(ashwin): Create auction ID -> name authority index.

		status = types.AuthorityUnderAuction
		auctionID = auction.ID
	}

	authority := types.NameAuthority{
		Height:         ctx.BlockHeight(),
		OwnerAddress:   owner.String(),
		OwnerPublicKey: helpers.BytesToBase64(pubKey.Bytes()),
		Status:         status,
		AuctionID:      auctionID,
	}

	k.SetNameAuthority(ctx, name, authority)

	return nil
}

// ProcessReserveSubAuthority reserves a sub-authority.
func (k Keeper) ProcessReserveSubAuthority(ctx sdk.Context, name string, msg types.MsgReserveAuthority) (string, sdk.Error) {
	// Get parent authority name.
	names := strings.Split(name, ".")
	parent := strings.Join(names[1:], ".")

	// Check if parent authority exists.
	parentAuthority := k.GetNameAuthority(ctx, parent)
	if parentAuthority == nil {
		return name, sdk.ErrInternal("Parent authority not found.")
	}

	// Sub-authority creator needs to be the owner of the parent authority.
	if parentAuthority.OwnerAddress != msg.Signer.String() {
		return name, sdk.ErrUnauthorized("Access denied.")
	}

	// Sub-authority owner defaults to parent authority owner.
	subAuthorityOwner := msg.Signer
	if !msg.Owner.Empty() {
		// Override sub-authority owner if provided in message.
		subAuthorityOwner = msg.Owner
	}

	sdkErr := k.createAuthority(ctx, name, subAuthorityOwner, false)
	if sdkErr != nil {
		return "", sdkErr
	}

	return name, nil
}

func (k Keeper) checkWRNAccess(ctx sdk.Context, signer sdk.AccAddress, inputWRN string) sdk.Error {
	parsedWRN, err := url.Parse(inputWRN)
	if err != nil {
		return sdk.ErrInternal("Invalid WRN.")
	}

	name := parsedWRN.Host
	formattedWRN := fmt.Sprintf("wrn://%s%s", name, parsedWRN.RequestURI())
	if formattedWRN != inputWRN {
		return sdk.ErrInternal("Invalid WRN.")
	}

	// Check authority record.
	authority := k.GetNameAuthority(ctx, name)
	if authority == nil {
		return sdk.ErrInternal("Name authority not found.")
	}

	if authority.OwnerAddress != signer.String() {
		return sdk.ErrUnauthorized("Access denied.")
	}

	if authority.Status != types.AuthorityActive {
		return sdk.ErrUnauthorized("Authority is not active.")
	}

	return nil
}

// ProcessSetName creates a WRN -> Record ID mapping.
func (k Keeper) ProcessSetName(ctx sdk.Context, msg types.MsgSetName) sdk.Error {
	err := k.checkWRNAccess(ctx, msg.Signer, msg.WRN)
	if err != nil {
		return err
	}

	nameRecord := k.GetNameRecord(ctx, msg.WRN)
	if nameRecord != nil && nameRecord.ID == msg.ID {
		// Already pointing to same ID, no-op.
		return nil
	}

	k.SetNameRecord(ctx, msg.WRN, msg.ID)

	return nil
}

// ProcessDeleteName removes a WRN -> Record ID mapping.
func (k Keeper) ProcessDeleteName(ctx sdk.Context, msg types.MsgDeleteName) sdk.Error {
	err := k.checkWRNAccess(ctx, msg.Signer, msg.WRN)
	if err != nil {
		return err
	}

	if !k.HasNameRecord(ctx, msg.WRN) {
		return sdk.ErrInternal("Name not found.")
	}

	// Set CID to empty string.
	k.SetNameRecord(ctx, msg.WRN, "")

	return nil
}
