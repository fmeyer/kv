package cmd

import (
	"fmt"
	"log"

	"github.com/fmeyer/kv/internal/db"

	"github.com/spf13/cobra"
)

var (
	key, value string
	kv         *db.KV
)

var rootCmd = &cobra.Command{Use: "kv"}

var setCmd = &cobra.Command{
	Use:   "s",
	Short: "Set a value",
	Run: func(cmd *cobra.Command, args []string) {
		kv.Set(key, value)
	},
}

var getCmd = &cobra.Command{
	Use:   "g",
	Short: "Get a value",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(kv.Get(key))
	},
}

var listCmd = &cobra.Command{
	Use:   "l",
	Short: "List all keys",
	Run: func(cmd *cobra.Command, args []string) {
		kv.List()
	},
}

func Execute() {
	defer kv.Close()

	setCmd.Flags().StringVarP(&key, "key", "k", "", "Key")
	setCmd.Flags().StringVarP(&value, "value", "v", "", "Value")
	setCmd.MarkFlagRequired("key")
	setCmd.MarkFlagRequired("value")

	getCmd.Flags().StringVarP(&key, "key", "k", "", "Key")
	getCmd.MarkFlagRequired("key")

	rootCmd.AddCommand(setCmd, getCmd, listCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	kv = db.NewKV()
}
