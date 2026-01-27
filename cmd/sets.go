package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/config"
	"github.com/spf13/cobra"
)

var setsCmd = &cobra.Command{
	Use:   "sets",
	Short: "Manage global sets",
	Long:  "Interactively sets the global sets priority list. If no --sets argument is used with apply, global will be used. If a card print is not found from the --sets argument list, global sets will be checked.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("a subcommand is required (add | remove | list)")
	},
}

var setsAddCmd = &cobra.Command{
	Use:   "add <set>",
	Short: "Add a set[s] to priority list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sets := parseCSV(args[0])
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		existing := make(map[string]struct{}, len(cfg.GlobalSetPriority))
		for _, existingSet := range cfg.GlobalSetPriority {
			existing[existingSet] = struct{}{}
		}

		setsToAdd := make([]string, 0)
		for _, set := range sets {
			tSet := strings.ToUpper(set)
			if _, ok := existing[tSet]; ok {
				continue
			}
			existing[tSet] = struct{}{}
			setsToAdd = append(setsToAdd, tSet)
			fmt.Printf("Adding %v\n", tSet)
		}

		if len(setsToAdd) == 0 {
			println("No sets to add")
			return nil
		}

		cfg.GlobalSetPriority = append(cfg.GlobalSetPriority, setsToAdd...)
		err = cfg.SaveConfig()
		if err != nil {
			return err
		}

		return nil
	},
}

var setsRemoveCmd = &cobra.Command{
	Use:   "remove <set>",
	Short: "Remove a set[s] from priority list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sets := parseCSV(args[0])
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		remove := make(map[string]struct{}, len(sets))
		for _, set := range sets {
			remove[strings.ToUpper(set)] = struct{}{}
		}

		kept := cfg.GlobalSetPriority[:0]
		for _, existing := range cfg.GlobalSetPriority {
			if _, ok := remove[existing]; !ok {
				kept = append(kept, existing)
			} else {
				fmt.Printf("Removing %v", existing)
			}
		}

		if len(kept) == len(cfg.GlobalSetPriority) {
			println("No sets to remove")
			return nil
		}

		cfg.GlobalSetPriority = kept
		err = cfg.SaveConfig()
		if err != nil {
			return err
		}

		return nil
	},
}

var setsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List global sets",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		if len(cfg.GlobalSetPriority) == 0 {
			println("No global sets have been added")
		}

		for _, set := range cfg.GlobalSetPriority {
			println(set)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setsCmd)
	setsCmd.AddCommand(setsAddCmd)
	setsCmd.AddCommand(setsRemoveCmd)
	setsCmd.AddCommand(setsListCmd)

	setsCmd.SilenceUsage = true
}
