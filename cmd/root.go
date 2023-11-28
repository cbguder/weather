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
	startDate *time.Time
	endDate   *time.Time
)

var rootCmd = &cobra.Command{
	Use: "weather",

	PersistentPreRunE: rootPreRunE,
}

func init() {
	rootCmd.PersistentFlags().String("after", "", "only use records after this date")
	rootCmd.PersistentFlags().String("before", "", "only use records before this date")
	rootCmd.PersistentFlags().String("cache", "", "cache directory (default ~/.cache/weather)")
	rootCmd.PersistentFlags().Bool("ignore-cache", false, "ignore cached data")
}

func rootPreRunE(cmd *cobra.Command, _ []string) error {
	noaa.IgnoreCache, _ = cmd.Flags().GetBool("ignore-cache")
	noaa.CacheDir = getCacheDir(cmd)

	startDate = parseDateFlag(cmd, "after")
	endDate = parseDateFlag(cmd, "before")

	if startDate != nil && endDate != nil && startDate.After(*endDate) {
		return errors.New("start date must be before end date")
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
