package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/goccy/go-yaml"
	"github.com/nabeken/psadm/v2"
)

type ExportCommand struct {
	KeyPrefix string `long:"key-prefix" description:"Specify a key prefix to be exported"`
}

func (cmd *ExportCommand) Execute(args []string) error {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := psadm.NewClient(cfg)

	params, err := client.GetParametersByPath(ctx, cmd.KeyPrefix)
	if err != nil {
		return err
	}

	out, err := yaml.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshaling into YAML: %w", err)
	}

	fmt.Print(string(out))

	return nil
}

func init() {
	parser.AddCommand(
		"export",
		"Export parameters",
		"The export command exports parameters from Parameter Store.",
		&ExportCommand{},
	)
}
