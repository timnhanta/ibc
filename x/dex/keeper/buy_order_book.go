package keeper

import (
	"interchange/x/dex/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetBuyOrderBook set a specific buyOrderBook in the store from its index
func (k Keeper) SetBuyOrderBook(ctx sdk.Context, buyOrderBook types.BuyOrderBook) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BuyOrderBookKeyPrefix))
	b := k.cdc.MustMarshal(&buyOrderBook)
	store.Set(types.BuyOrderBookKey(
		buyOrderBook.Index,
	), b)
}

// GetBuyOrderBook returns a buyOrderBook from its index
func (k Keeper) GetBuyOrderBook(
	ctx sdk.Context,
	index string,

) (val types.BuyOrderBook, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BuyOrderBookKeyPrefix))

	b := store.Get(types.BuyOrderBookKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveBuyOrderBook removes a buyOrderBook from the store
func (k Keeper) RemoveBuyOrderBook(
	ctx sdk.Context,
	index string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BuyOrderBookKeyPrefix))
	store.Delete(types.BuyOrderBookKey(
		index,
	))
}

// GetAllBuyOrderBook returns all buyOrderBook
func (k Keeper) GetAllBuyOrderBook(ctx sdk.Context) (list []types.BuyOrderBook) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BuyOrderBookKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.BuyOrderBook
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (b *BuyOrderBook) LiquidateFromSellOrder(order Order) (
	remainingSellOrder Order,
	liquidatedBuyOrder Order,
	gain int32,
	match bool,
	filled bool,
) {
	remainingSellOrder = order

	// No match if no order
	orderCount := len(b.Book.Orders)
	if orderCount == 0 {
		return order, liquidatedBuyOrder, gain, false, false
	}

	// Check if match
	highestBid := b.Book.Orders[orderCount-1]
	if order.Price > highestBid.Price {
		return order, liquidatedBuyOrder, gain, false, false
	}

	liquidatedBuyOrder = *highestBid

	// Check if sell order can be entirely filled
	if highestBid.Amount >= order.Amount {
		remainingSellOrder.Amount = 0
		liquidatedBuyOrder.Amount = order.Amount
		gain = order.Amount * highestBid.Price

		// Remove the highest bid if it has been entirely liquidated
		highestBid.Amount -= order.Amount
		if highestBid.Amount == 0 {
			b.Book.Orders = b.Book.Orders[:orderCount-1]
		} else {
			b.Book.Orders[orderCount-1] = highestBid
		}

		return remainingSellOrder, liquidatedBuyOrder, gain, true, true
	}

	// Not entirely filled
	gain = highestBid.Amount * highestBid.Price
	b.Book.Orders = b.Book.Orders[:orderCount-1]
	remainingSellOrder.Amount -= highestBid.Amount

	return remainingSellOrder, liquidatedBuyOrder, gain, true, false
}
