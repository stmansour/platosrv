package main

// Tickers is a map of the stock tickers with data loaded into the database.
// We could have scanned the database to find this information, but that would
// take a lot of time.  This makes things much faster.
// -----------------------------------------------------------------------------
var Tickers = map[string]int{
	"EURUSD": 0,
	"GBPUSD": 0,
	"USDCHF": 0,
	"USDJPY": 0,
	"EURGBP": 0,
	"EURCHF": 0,
	"EURJPY": 0,
	"GBPCHF": 0,
	"GBPJPY": 0,
	"CHFJPY": 0,
	"USDCAD": 0,
	"EURCAD": 0,
	"AUDUSD": 1, // ON
	"AUDJPY": 0,
	"NZDUSD": 0,
	"NZDJPY": 1, // ON
	"XAUUSD": 0,
	"XAGUSD": 0,
	"USDCZK": 0,
	"USDDKK": 0,
	"EURRUB": 0,
	"USDHUF": 0,
	"USDNOK": 0,
	"USDPLN": 0,
	"USDRUB": 0,
	"USDSEK": 0,
	"USDSGD": 0,
	"USDZAR": 0,
	"USDTRY": 0,
	"EURTRY": 0,
	"EURAUD": 0,
	"EURNZD": 0,
	"EURSGD": 0,
	"EURZAR": 0,
	"XAUEUR": 0,
	"XAGEUR": 0,
	"GBPCAD": 0,
	"GBPAUD": 1, // ON
	"GBPNZD": 0,
	"AUDCHF": 0,
	"AUDCAD": 0,
	"AUDNZD": 0,
	"NZDCHF": 0,
	"NZDCAD": 0,
	"CADCHF": 0,
	"CADJPY": 0,
	"USDUAH": 0,
	"WTIUSD": 0,
	"DJIUSD": 0,
	"SPXUSD": 0,
	"NDQUSD": 0,
	"USXUSD": 0,
	"USDHKD": 0,
	"EURHKD": 0,
	"USDMXN": 0,
	"EURMXN": 0,
	"USDILS": 0,
	"EURILS": 0,
	"BTCUSD": 0,
	"BRNUSD": 0,
	"USDCNH": 0,
	"GASUSD": 0,
	"EURNOK": 0,
	"EURSEK": 0,
	"EURDKK": 0,
	"EURCZK": 0,
	"EURHUF": 0,
	"EURPLN": 0,
	"ETHUSD": 0,
	"LTCUSD": 0,
	"EURUAH": 0,
	"EURCNH": 0,
	"BTCEUR": 0,
	"ETHEUR": 0,
	"LTCEUR": 0,
	"USDWMR": 0,
	"USDWMU": 0,
}
