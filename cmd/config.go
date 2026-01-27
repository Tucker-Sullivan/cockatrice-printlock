package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/config"
	"github.com/spf13/cobra"
)

var configSetters = map[string]func(*config.Config, string) error{
	"deckDir": func(c *config.Config, input string) error {
		c.DeckDir = input
		return nil
	},
	"cardsXmlPath": func(c *config.Config, input string) error {
		c.CardsXMLPath = input
		return nil
	},
	"globalSetPriority": func(c *config.Config, input string) error {
		parsed := parseCSV(strings.ToUpper(input))
		c.GlobalSetPriority = parsed
		return nil
	},
	"globalSetsHavePriority": func(c *config.Config, input string) error {
		value, err := checkBoolInput(input)
		if err != nil {
			return err
		}
		c.GlobalSetsHavePriority = value
		return nil
	},
	"preferHigherCollectionNumber": func(c *config.Config, input string) error {
		value, err := checkBoolInput(input)
		if err != nil {
			return err
		}
		c.PreferHigherCollectionNumber = value
		return nil
	},
}

func checkBoolInput(input string) (bool, error) {
	input = strings.TrimSpace(strings.ToLower(input))
	value, ableToParse := false, true
	switch input {
	case "true":
		value = true
	case "false":
		value = false
	default:
		ableToParse = false
	}
	if !ableToParse {
		return false, errors.New("value must be true or false")
	}
	return value, nil
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage config file",
	Long:  "List or set values from the config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("a subcommand is required (set | list)")
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List current config",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(cfg, "", " ")
		if err != nil {
			return err
		}

		println((string)(b))

		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set value for key",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		setter, ok := configSetters[args[0]]
		if !ok {
			return fmt.Errorf("unknown config key: %v", args[0])
		}

		if err := setter(cfg, args[1]); err != nil {
			return err
		}

		if err := cfg.SaveConfig(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configSetCmd)

	configCmd.SilenceUsage = true
}
