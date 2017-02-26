package main

import (
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/nabeken/psadm/ps"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type ImportCommand struct {
	Dryrun          bool   `long:"dryrun" description:"Perform dryrun"`
	Overwrite       bool   `long:"overwrite" description:"Overwrite the value in the key if it exists"`
	SkipExist       bool   `long:"skip-exist" description:"Skip the existing key"`
	DefaultKMSKeyID string `long:"default-kms-key-id" description:"Specify a default KMS Key ID"`
}

func (cmd *ImportCommand) Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("You must specify a YAML file to be imported.")
	}

	f, err := os.Open(args[0])
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", args[0])
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.Wrapf(err, "failed to read data from %s", args[0])
	}

	var params []*ps.Parameter
	if err := yaml.Unmarshal(data, &params); err != nil {
		return errors.Wrap(err, "failed to unmarshal from YAML")
	}

	client := ps.NewClient(session.Must(session.NewSession()))
	for _, p := range params {
		if p.Type == ssm.ParameterTypeSecureString && p.KMSKeyID == "" {
			p.KMSKeyID = cmd.DefaultKMSKeyID
		}
		if err := client.PutParameter(p, cmd.Overwrite); err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				if awsErr.Code() == ssm.ErrCodeParameterAlreadyExists && cmd.SkipExist {
					continue
				}
			}
			return errors.Wrapf(err, "failed to update '%s'", p.Name)
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
