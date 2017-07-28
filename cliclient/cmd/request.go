package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aakso/ssh-inscribe/pkg/client"
	"github.com/spf13/cobra"
)

var ReqCmd = &cobra.Command{
	Use:   "req",
	Short: "Login to server and generate SSH certificate",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := &client.Client{
			Config: ClientConfig,
		}
		defer c.Close()
		if b, _ := cmd.Flags().GetBool("clear"); b == true {
			return c.Logout()
		} else if b, _ := cmd.Flags().GetBool("list-logins"); b == true {
			discoverResult, err := c.GetAuthenticators()
			if err != nil {
				return err
			}
			for _, v := range discoverResult {
				fmt.Printf("%s (%s)\n", v.AuthenticatorName, v.AuthenticatorRealm)
			}
			return nil
		}
		return c.Login()
	},
}

func init() {
	RootCmd.AddCommand(ReqCmd)
	ReqCmd.Flags().StringVarP(
		&ClientConfig.IdentityFile,
		"identity",
		"i",
		os.Getenv("SSH_INSCRIBE_IDENTITY"),
		"Identity (private key) file location. Required if --generate is unset ($SSH_INSCRIBE_IDENTITY)",
	)

	var defExpire time.Duration
	if expire := os.Getenv("SSH_INSCRIBE_EXPIRE"); expire != "" {
		defExpire, _ = time.ParseDuration(expire)
	}
	ReqCmd.Flags().DurationVarP(
		&ClientConfig.CertLifetime,
		"expire",
		"e",
		defExpire,
		"Request specific lifetime. Example '10m' ($SSH_INSCRIBE_EXPIRE)",
	)
	if os.Getenv("SSH_INSCRIBE_WRITE") != "" {
		ClientConfig.WriteCert = true
	}
	ReqCmd.Flags().BoolVarP(
		&ClientConfig.WriteCert,
		"write",
		"w",
		ClientConfig.WriteCert,
		"Write certificate (and generated keys) to file specified by <identity> ($SSH_INSCRIBE_WRITE)",
	)

	if os.Getenv("SSH_INSCRIBE_RENEW") != "" {
		ClientConfig.AlwaysRenew = true
	}
	ReqCmd.Flags().BoolVar(
		&ClientConfig.AlwaysRenew,
		"renew",
		ClientConfig.AlwaysRenew,
		"Always renew the certificate even if it is not expired ($SSH_INSCRIBE_RENEW)",
	)

	if os.Getenv("SSH_INSCRIBE_USE_AGENT") == "0" {
		ClientConfig.UseAgent = false
	}
	ReqCmd.Flags().BoolVar(
		&ClientConfig.UseAgent,
		"agent",
		ClientConfig.UseAgent,
		"Store key and certificate to a ssh-agent specified by $SSH_AUTH_SOCK ($SSH_INSCRIBE_USE_AGENT)",
	)

	if os.Getenv("SSH_INSCRIBE_GENKEY") != "" {
		ClientConfig.GenerateKeypair = true
	}
	ReqCmd.Flags().BoolVarP(
		&ClientConfig.GenerateKeypair,
		"generate",
		"g",
		ClientConfig.GenerateKeypair,
		"Generate ad-hoc keypair. Useful with ssh-agent ($SSH_INSCRIBE_GENKEY)",
	)

	ReqCmd.Flags().Bool(
		"clear",
		false,
		"Clear granted certificate",
	)

	ReqCmd.Flags().Bool(
		"list-logins",
		false,
		"List available auth endpoints",
	)

	defLoginAuthEndpoints := []string{}
	if logins := os.Getenv("SSH_INSCRIBE_LOGIN_AUTH_ENDPOINTS"); logins != "" {
		defLoginAuthEndpoints = strings.Split(logins, ",")
	}
	ReqCmd.Flags().StringSliceVarP(
		&ClientConfig.LoginAuthEndpoints,
		"login",
		"l",
		defLoginAuthEndpoints,
		"Login to specific auth endpoits ($SSH_INSCRIBE_LOGIN_AUTH_ENDPOINTS)",
	)

	var defIncludePrincipals string
	if s := os.Getenv("SSH_INSCRIBE_INCLUDE_PRINCIPALS"); s != "" {
		defIncludePrincipals = s
	}
	ReqCmd.Flags().StringVar(
		&ClientConfig.IncludePrincipals,
		"include",
		defIncludePrincipals,
		"Request only principals matching the glob pattern to be included ($SSH_INSCRIBE_INCLUDE_PRINCIPALS)",
	)

	var defExcludePrincipals string
	if s := os.Getenv("SSH_INSCRIBE_EXCLUDE_PRINCIPALS"); s != "" {
		defExcludePrincipals = s
	}
	ReqCmd.Flags().StringVar(
		&ClientConfig.ExcludePrincipals,
		"exclude",
		defExcludePrincipals,
		"Request only principals not matching the glob pattern to be included ($SSH_INSCRIBE_EXCLUDE_PRINCIPALS)",
	)
}
