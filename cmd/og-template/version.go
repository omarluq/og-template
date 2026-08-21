package main

import (
	"fmt"

	"github.com/samber/oops"
	"github.com/spf13/cobra"

	"github.com/omarluq/og-template/internal/vinfo"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print build version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), vinfo.String()); err != nil {
				return oops.Wrapf(err, "write version output")
			}

			return nil
		},
	}
}
