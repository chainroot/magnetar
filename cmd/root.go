package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	homePath     string
	dataDir      string
	cosmosSdk    bool
	cometBft     bool
	txIndex      bool
	keepBlocks   uint64
	keepVersions uint64
	compact      bool

	appName = "magnetar"
)

func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   appName,
		Short: "magnetar cleans up databases of Cosmos SDK applications, removing historical data generally not needed for validator nodes",
	}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		// reads `homeDir/config.yaml` into `var config *Config` before each command
		// if err := initConfig(rootCmd); err != nil {
		// 	return err
		// }

		return nil
	}

	// --keep-blocks flag
	rootCmd.PersistentFlags().Uint64VarP(&keepBlocks, "keep-blocks", "b", 10, "set the amount of blocks to keep")
	if err := viper.BindPFlag("keep-blocks", rootCmd.PersistentFlags().Lookup("keep-blocks")); err != nil {
		panic(err)
	}

	// --keep-versions flag
	rootCmd.PersistentFlags().Uint64VarP(&keepVersions, "keep-versions", "v", 10, "set the amount of versions to keep in the application store")
	if err := viper.BindPFlag("keep-versions", rootCmd.PersistentFlags().Lookup("keep-versions")); err != nil {
		panic(err)
	}

	// --cosmos-sdk flag
	rootCmd.PersistentFlags().BoolVar(&cosmosSdk, "cosmos-sdk", true, "set to false if using only with cometbft")
	if err := viper.BindPFlag("cosmos-sdk", rootCmd.PersistentFlags().Lookup("cosmos-sdk")); err != nil {
		panic(err)
	}

	// --cometbft flag
	rootCmd.PersistentFlags().BoolVar(&cometBft, "cometbft", true, "set to false you dont want to prune cometbft data")
	if err := viper.BindPFlag("cometbft", rootCmd.PersistentFlags().Lookup("cometbft")); err != nil {
		panic(err)
	}

	// --tx-index flag
	rootCmd.PersistentFlags().BoolVar(&txIndex, "tx-index", true, "set to false you dont want to prune tx-index data")
	if err := viper.BindPFlag("tx-index", rootCmd.PersistentFlags().Lookup("tx-index")); err != nil {
		panic(err)
	}

	// --compact flag
	rootCmd.PersistentFlags().BoolVar(&compact, "compact", false, "set to true if you want to compact all databases")
	if err := viper.BindPFlag("compact", rootCmd.PersistentFlags().Lookup("compact")); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(pruneCmd())

	return rootCmd
}

func Execute() {
	cobra.EnableCommandSorting = false

	rootCmd := NewRootCmd()
	rootCmd.SilenceUsage = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
