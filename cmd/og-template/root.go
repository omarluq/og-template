package main

import "github.com/spf13/cobra"

func newRootCmd() *cobra.Command {
	var cfgFile string

	cmd := &cobra.Command{
		Use:           "og-template",
		Short:         "og-template is a command line tool",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
	cmd.AddCommand(newConfigCmd(&cfgFile))
	cmd.AddCommand(newVersionCmd())

	return cmd
}
