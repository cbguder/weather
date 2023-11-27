package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/cbguder/weather/noaa"
)

var (
	afterDate  *time.Time
	beforeDate *time.Time
)

var rootCmd = &cobra.Command{
	Use: "weather",

	PersistentPreRunE: rootPreRunE,
}

func init() {
	rootCmd.PersistentFlags().String("after", "", "only use records after this date")
	rootCmd.PersistentFlags().String("before", "", "only use records before this date")
	rootCmd.PersistentFlags().String("cache", "", "cache directory (default ~/.cache/weather)")
}

func rootPreRunE(cmd *cobra.Command, _ []string) error {
	noaa.CacheDir = getCacheDir(cmd)

	afterDate = parseDateFlag(cmd, "after")
	beforeDate = parseDateFlag(cmd, "before")

	if afterDate != nil && beforeDate != nil && beforeDate.Before(*afterDate) {
		return errors.New("before date must be after after date")
	}

	return nil
}

func parseDateFlag(cmd *cobra.Command, name string) *time.Time {
	date, _ := cmd.Flags().GetString(name)
	if date == "" {
		return nil
	}

	t, err := time.Parse(time.DateOnly, date)
	cobra.CheckErr(err)
	return &t
}

func getCacheDir(cmd *cobra.Command) string {
	cacheDir, _ := cmd.Flags().GetString("cache")
	if cacheDir != "" {
		return cacheDir
	}

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	return filepath.Join(home, ".cache", "weather")
}

func Execute() error {
	return rootCmd.Execute()
}
