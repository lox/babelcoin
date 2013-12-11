package main

import (
	"github.com/spf13/cobra"
	"github.com/davecgh/go-spew/spew"
	"./btce"
	"os"
)

func btceFromEnv() *btce.BtceApi {
	key := os.Getenv("BTCE_KEY")
	secret := os.Getenv("BTCE_SECRET")
	return btce.NewBtceApi(btce.PrivateApiUrl, key, secret)
}

func main() {
	var BtceCmd = &cobra.Command{
		Use:   "btce",
		Short: "Interface to btc-e.com",
		Long: "Interface to btc-e.com",
	}

	var Btce_BalancesCmd = &cobra.Command{
		Use:   "balances",
		Short: "Account balances on btc-e.com",
		Run: func(cmd *cobra.Command, args []string) {
			b := btceFromEnv()
			info, error := b.GetInfo();
			if error != nil {
				panic(error)
			}

			spew.Dump(info)
		},
	}

	var Btce_TransactionsCmd = &cobra.Command{
		Use:   "transactions",
		Short: "Transaction history on btc-e.com",
		Run: func(cmd *cobra.Command, args []string) {
			b := btceFromEnv()
			info, error := b.TransHistory(map[string]string{
				"count": "10",
			});
			if error != nil {
				panic(error)
			}

			spew.Dump(info)
		},
	}

	var Btce_TradesCmd = &cobra.Command{
		Use:   "trades",
		Short: "Trade history on btc-e.com",
		Run: func(cmd *cobra.Command, args []string) {
			b := btceFromEnv()
			info, error := b.TradeHistory(map[string]string{
				"count": "10",
			});
			if error != nil {
				panic(error)
			}

			spew.Dump(info)
		},
	}

	var Btce_TradeCmd = &cobra.Command{
		Use:   "trade",
		Short: "Conduct a round-trip trade on btc-e.com",
		Run: func(cmd *cobra.Command, args []string) {
			b := btceFromEnv()
			resp, error := b.Trade("ftc_btc", "buy", 0.0001, 1)
			if error != nil {
				panic(error)
			}
			spew.Dump(resp)

			resp2, error2 := b.CancelOrder(resp.OrderId)
			if error2 != nil {
				panic(error2)
			}
			spew.Dump(resp2)
		},
	}

	var Btce_TickerCmd = &cobra.Command{
		Use:   "ticker",
		Short: "Poll the ticker on btc-e.com",
		Run: func(cmd *cobra.Command, args []string) {
			//i := btce.NewInfoApi(btce.InfoUrl)
			//info, _ := i.Pairs()
			//spew.Dump(info)

			t := btce.NewTradesApi(btce.TradesUrl, []string{"ltc_btc"}, 10)
			trades, _ := t.Trades()

			spew.Dump(trades)

			b := btce.NewTickerApi(btce.TickerUrl, []string{"ltc_btc"})
			candles, error := b.Candles()
			if error != nil {
				panic(error)
			}

			spew.Dump(candles)
		},
	}

	BtceCmd.AddCommand(Btce_BalancesCmd)
	BtceCmd.AddCommand(Btce_TransactionsCmd)
	BtceCmd.AddCommand(Btce_TradesCmd)
	BtceCmd.AddCommand(Btce_TradeCmd)
	BtceCmd.AddCommand(Btce_TickerCmd)

	var rootCmd = &cobra.Command{Use: "babelcoin"}
	rootCmd.AddCommand(BtceCmd)
	rootCmd.Execute()
}