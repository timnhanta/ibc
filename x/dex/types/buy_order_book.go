// x/dex/types/buy_order_book.go

package types

func NewBuyOrderBook(AmountDenom string, PriceDenom string) BuyOrderBook {
	book := NewOrderBook()
	return BuyOrderBook{
		AmountDenom: AmountDenom,
		PriceDenom:  PriceDenom,
		Book:        &book,
	}
}

func (b *BuyOrderBook) AppendOrder(creator string, amount int32, price int32) (int32, error) {
	return b.Book.appendOrder(creator, amount, price, Increasing)
}
