package cmd

import (
	"errors"
	"time"

	"github.com/spf13/cobra"
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
}

func rootPreRunE(cmd *cobra.Command, _ []string) error {
	var err error

	afterDate, err = parseDateFlag(cmd, "after")
	if err != nil {
		return err
	}

	beforeDate, err = parseDateFlag(cmd, "before")
	if err != nil {
		return err
	}

	if afterDate != nil && beforeDate != nil && beforeDate.Before(*afterDate) {
		return errors.New("before date must be after after date")
	}

	return nil
}

func parseDateFlag(cmd *cobra.Command, name string) (*time.Time, error) {
	date, _ := cmd.Flags().GetString(name)
	if date == "" {
		return nil, nil
	}

	t, err := time.Parse(time.DateOnly, date)
	return &t, err
}

func Execute() error {
	return rootCmd.Execute()
}
