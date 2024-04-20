package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/goccy/go-yaml"
	"github.com/nabeken/psadm/v2"
)

type ImportCommand struct {
	Dryrun          bool   `long:"dryrun" description:"Perform dryrun"`
	Overwrite       bool   `long:"overwrite" description:"Overwrite the value in the key if it exists"`
	SkipExist       bool   `long:"skip-exist" description:"Skip the existing key"`
	DefaultKMSKeyID string `long:"default-kms-key-id" description:"Specify a default KMS Key ID"`
}

func (cmd *ImportCommand) Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("You must specify a YAML file to be imported")
	}

	f, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("opening %s: %w", args[0], err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("reading from %s: %w", args[0], err)
	}

	var params []*psadm.Parameter
	if err := yaml.Unmarshal(data, &params); err != nil {
		return fmt.Errorf("marshaling into YAML: %w", err)
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := psadm.NewClient(cfg)

	// function to update
	actualRun := func(p *psadm.Parameter) error {
		if err := client.PutParameter(ctx, p, cmd.Overwrite); err != nil {
			var ae *types.ParameterAlreadyExists

			if errors.As(err, &ae) && cmd.SkipExist {
				return nil
			}

			return err
		}
		return nil
	}
	dryRun := func(p *psadm.Parameter) error {
		fmt.Printf("dryrun: '%s' will be updated\n", p.Name)
		return nil
	}

	runF := actualRun
	if cmd.Dryrun {
		runF = dryRun
	}

	for _, p := range params {
		if p.Type == string(types.ParameterTypeSecureString) && p.KMSKeyID == "" {
			p.KMSKeyID = cmd.DefaultKMSKeyID
		}
		if err := runF(p); err != nil {
			return fmt.Errorf("updating '%s': %w", p.Name, err)
		}
	}

	return nil
}

func init() {
	parser.AddCommand(
		"import",
		"Import parameters",
		"The import command imports parameters from exported YAML file.",
		&ImportCommand{},
	)
}
