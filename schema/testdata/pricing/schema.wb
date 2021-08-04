// Candlestick
struct Candle {
	open 	f64		// Open price
	high 	f64		// High price
	low 	f64		// Low price
	close 	f64		// Close price
}

struct Spread {
	low 	f64		// Low spread
	high 	f64		// High spread
	avg 	f64		// Avg spread
}

struct Liquidations {
	trades		i64		// Trades
	min			f64		// Min price
	avg			f64		// Avg price
	max			f64		// Max price
	buys		f64		// Buys
	sells		f64		// Sells
	value		f64		// Value
}

// Greeks are financial measures of the sensitivity of an option’s price to its
// underlying determining parameters, such as volatility or the price of the underlying
// asset. The Greeks are utilized in the analysis of an options portfolio and in sensitivity
// analysis of an option or portfolio of options. The measures are considered essential by
// many investors for making informed decisions in options trading.
//
// Delta, Gamma, Vega, Theta, and Rho are the key option Greeks. However, there are many other
// option Greeks that can be derived from those mentioned above.
struct Greeks {
	// Implied Volatility
	iv		f64
	// Delta (Δ) is a measure of the sensitivity of an option’s price changes relative to the
	// changes in the underlying asset’s price. In other words, if the price of the underlying
	// asset increases by $1, the price of the option will change by Δ amount.
	delta	f64
	// Gamma (Γ) is a measure of the delta’s change relative to the changes in the price of the
	// underlying asset. If the price of the underlying asset increases by $1, the option’s delta
	// will change by the gamma amount. The main application of gamma is the assessment of the
	// option’s delta.
	gamma	f64
	// Vega (ν) is an option Greek that measures the sensitivity of an option price relative to
	// the volatility of the underlying asset. If the volatility of the underlying asses increases
	// by 1%, the option price will change by the vega amount.
	vega	f64
	// Theta (θ) is a measure of the sensitivity of the option price relative to the option’s time
	// to maturity. If the option’s time to maturity decreases by one day, the option’s price will
	// change by the theta amount. The Theta option Greek is also referred to as time decay.
	theta	f64
	// Rho (ρ) measures the sensitivity of the option price relative to interest rates. If a benchmark
	// interest rate increases by 1%, the option price will change by the rho amount. The rho is
	// considered the least significant among other option Greeks because option prices are generally
	// less sensitive to interest rate changes than to changes in other parameters.
	rho		f64
}

struct Ticks {
	total	i64
	up		i64
	down	i64
}

struct Volume {
	total	f64
	buy		f64
	sell	f64
}

struct Time {
	start 		i64		// Start timestamp (UTC unix millis)
	end 		i64		// End timestamp (UTC unix millis)
}

struct Trades {
	count		i64		// Number of trades
	min			i64		// Min-trade ID (broker specific)
	max			i64		// Max-trade ID (broker specific)
}

enum Aggressor : i32 {
	Buy = 0
	Sell = 1
}

struct Trade {
	id			i64
	price		f64
	quantity	f64
	aggressor 	Aggressor
}

struct Bar {
	start 			i64			// Start timestamp (UTC unix millis)
    end 			i64			// End timestamp (UTC unix millis)
	price 			Candle		// Price candle / Mid-point
	bid 			Candle		// Bid candle
	ask 			Candle		// Ask candle
	spread 			Spread		// Spread
	ticks 			Ticks		// Ticks
	volume 			Volume		// Volume
	interest		Volume		// Open interest
}

struct OptionBar {
	start 			i64			// Start timestamp (UTC unix millis)
    end 			i64			// End timestamp (UTC unix millis)
	price 			Candle		// Price candle / Mid-point
	bid 			Candle		// Bid candle
	ask 			Candle		// Ask candle
	spread 			Spread		// Spread
	ticks 			Ticks		// Ticks
	volume 			Volume		// Volume
	interest		Volume		// Open interest
	greeks			Greeks		// Greeks
}
