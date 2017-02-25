package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/nabeken/psadm/ps"
	"gopkg.in/yaml.v2"
)

type ExportCommand struct {
	WithDecryption bool   `long:"with-decryption" description:"Decrypt a SecureString"`
	KeyPrefix      string `long:"key-prefix" description:"Specify a key prefix to be exported"`
}

func (cmd *ExportCommand) Execute(args []string) error {
	client := &ps.Client{ssm.New(session.Must(session.NewSession()))}

	params, err := client.GetAllParameters(cmd.KeyPrefix, cmd.WithDecryption)
	if err != nil {
		return err
	}

	out, err := yaml.Marshal(params)
	if err != nil {
		return err
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
