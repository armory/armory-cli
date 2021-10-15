package accounts

import (
	"context"
	"fmt"
	de "github.com/armory-io/deploy-engine/pkg"
	"github.com/spf13/cobra"
	_nethttp "net/http"
	"time"
)

const (
	listShort   = ""
	listLong    = ""
	listExample    = ""
)

type listOptions struct {
	*accountsOptions
}

func NewListCmd(accountsOptions *accountsOptions) *cobra.Command {
	options := &listOptions{
		accountsOptions: accountsOptions,
	}
	command := &cobra.Command{
		Use:     "list —-provider [<provider name>] ",
		Aliases: []string{"list"},
		Short:   listShort,
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(cmd, options, args)
		},
	}
	command.Flags().StringVarP(&options.providerName, "provider", "p", "", "provider name")

	return command
}

type FormattableAccountsStartResponse struct {
	// The deployment's ID.
	Accounts []de.AccountV2Account `json:"accounts,omitempty" yaml:"accounts,omitempty"`
	httpResponse *_nethttp.Response
	err error
}

func newAccountsListResponse(raw *de.AccountV2ListAccountsResponse, response *_nethttp.Response,
	err error) FormattableAccountsStartResponse {
	deployment := FormattableAccountsStartResponse{
		Accounts: raw.GetAccounts(),
		httpResponse: response,
		err: err,
	}
	return deployment
}
func (u FormattableAccountsStartResponse) Get() interface{} {
	return u
}

func (u FormattableAccountsStartResponse) GetHttpResponse() *_nethttp.Response {
	return u.httpResponse
}

func (u FormattableAccountsStartResponse) GetFetchError() error {
	return u.err
}

func (u FormattableAccountsStartResponse) String() string {
	return fmt.Sprintf("%v", u.Accounts)
}


func list(cmd *cobra.Command, options *listOptions, args []string) error {
	ctx, cancel := context.WithTimeout(options.DeployClient.Context, time.Second * 5)
	defer cancel()

	// prepare request
	request := options.DeployClient.DeploymentServiceApi.DeploymentServiceListAccounts(ctx).Provider(options.providerName)
	// execute request
	raw, response, err := request.Execute()
	// create response object
	accounts := newAccountsListResponse(&raw, response, err)
	dataFormat, err := options.Output.Formatter(accounts)
	if err != nil {
		return fmt.Errorf("error trying to parse response: %s", err)
	}
	_, err = fmt.Fprintln(cmd.OutOrStdout(), dataFormat)
	return err
}