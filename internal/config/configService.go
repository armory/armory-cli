package config

import (
	"fmt"
	"github.com/armory-io/go-cloud-service/pkg/client"
	tls2 "github.com/armory-io/go-cloud-service/pkg/tls"
	"github.com/armory-io/go-cloud-service/pkg/token"
	"github.com/juju/ansiterm"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

const (
	ParamAddClientId       = "client-id"
	ParamAddSecret         = "secret"
	ParamAddTokenIssuerUrl = "token-issuer"
	ParamAddAudience       = "audience"
	ParamProvider = "provider"
	deployCliEnvVarPrefix = "ARMORY"
	deployCliConfigDir    = ".armory"
	deployCliConfigFile   = "config.yaml"
)

func Add(cmd *cobra.Command) error {
	n, err := cmd.Flags().GetString(ParamContext)
	if err != nil {
		return err
	}

	c, err := LoadConfig(false)
	if err != nil {
		return err
	}

	clientId, err := cmd.Flags().GetString(ParamAddClientId)
	if err != nil {
		return err
	}

	secret, err := cmd.Flags().GetString(ParamAddSecret)
	if err != nil {
		return err
	}

	tokenIssuer, err := cmd.Flags().GetString(ParamAddTokenIssuerUrl)
	if err != nil {
		return err
	}

	audience, err := cmd.Flags().GetString(ParamAddAudience)
	if err != nil {
		return err
	}

	ctx := Context{
		Name: n,
		Identity: token.Identity{
			Armory: token.ArmoryCloud{
				ClientId:       clientId,
				Secret:         secret,
				Audience:       audience,
				TokenIssuerUrl: tokenIssuer,
			},
		},
		Connection: client.Service{},
	}
	if ctx.Connection.Grpc, err = cmd.Flags().GetString(ParamEndpoint); err != nil {
		return err
	}
	if ctx.Connection.Insecure, err = cmd.Flags().GetBool(ParamPlaintext); err != nil {
		return err
	}
	if !ctx.Connection.Insecure {
		tls := &tls2.Settings{}
		ctx.Connection.Tls = tls
		if tls.InsecureSkipVerify, err = cmd.Flags().GetBool(ParamInsecure); err != nil {
			return err
		}
		if tls.ClientCertFile, err = cmd.Flags().GetString(ParamCert); err != nil {
			return err
		}
		if tls.CAcertFile, err = cmd.Flags().GetString(ParamCacert); err != nil {
			return err
		}
		if tls.ClientKeyFile, err = cmd.Flags().GetString(ParamKey); err != nil {
			return err
		}
		if tls.ClientKeyPassword, err = cmd.Flags().GetString(ParamKeyPassword); err != nil {
			return err
		}
	}

	if ctx.Connection.NoProxy, err = cmd.Flags().GetBool(ParamNoProxy); err != nil {
		return err
	}

	found := false
	for i := range c.Contexts {
		if c.Contexts[i].Name == n {
			c.Contexts[i] = ctx
			found = true
		}
	}
	if !found {
		c.Contexts = append(c.Contexts, ctx)
	}
	return SaveConfig(c)
}

type Config struct {
	Contexts       []Context `yaml:"contexts"`
	CurrentContext string    `yaml:"current-context"`
}

func (c *Config) getContext(name string) *Context {
	for _, ctx := range c.Contexts {
		if ctx.Name == name {
			return &ctx
		}
	}
	return nil
}

func (c *Config) checkCurrentContext() {
	if c.getContext(c.CurrentContext) == nil {
		if len(c.Contexts) > 0 {
			c.CurrentContext = c.Contexts[0].Name
		}
	}
}

type Context struct {
	Identity   token.Identity `yaml:"identity"`
	Name       string         `yaml:"name"`
	Connection client.Service `yaml:"connection"`
}

func (c *Context) NewConnection() client.Connection {
	return client.New(c.Connection, &c.Identity, log.StandardLogger())
}

func getHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	home = path.Join(home, deployCliConfigDir)
	return home
}

func getConfigFile() string {
	if cfg := os.Getenv(fmt.Sprintf("%s_CONFIG", deployCliEnvVarPrefix)); cfg != "" {
		return cfg
	}
	return path.Join(getHome(), deployCliConfigFile)
}

// LoadConfig loads configuration from the default or overridden location
// If withDefaults is specified, it applies default settings to each loaded
// context.
func LoadConfig(withDefaults bool) (*Config, error) {
	f := getConfigFile()
	fl, err := os.Open(f)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	buf, err := ioutil.ReadAll(fl)
	if err != nil {
		_ = os.Remove(f)
		return nil, err
	}

	c := &Config{}
	if withDefaults {
		// We deserialize to a map
		m := make(map[string]interface{})
		if err := yaml.Unmarshal(buf, m); err != nil {
			return nil, err
		}

		// First decode non context attributes
		if err := mapstructure.Decode(m, c); err != nil {
			return nil, err
		}

		// Reset - we're rebuilding the slice with defaults
		c.Contexts = nil
		ctxs, ok := m["contexts"].([]interface{})
		if ok {
			for i := range ctxs {
				ctx := Context{
					Connection: defaultService(),
					Identity:   token.DefaultIdentity(),
				}
				if err := mapstructure.Decode(ctxs[i], &ctx); err != nil {
					return nil, err
				}
				c.Contexts = append(c.Contexts, ctx)
			}
		}

	} else {
		if err := yaml.Unmarshal(buf, c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// SaveConfig saves the current configuration. It also checks the current context
// refers to an existing context.
func SaveConfig(c *Config) error {
	f := getConfigFile()
	c.checkCurrentContext()
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f, b, 0600)
}

// defaultService returns default connection settings
func defaultService() client.Service {
	return client.Service{
		Grpc:                      "deploy.cloud.armory.io:443",
		KeepAliveHeartbeatSeconds: 30,
		KeepAliveTimeOutSeconds:   10,
	}
}

func View() {
	c, err := LoadConfig(true)
	if err != nil {
		fmt.Println("Unable to read configuration file:", err.Error())
		return
	}

	if len(c.Contexts) == 0 {
		fmt.Println("Non existent or empty configuration file.")
		return
	}

	wt := ansiterm.NewTabWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprint(wt, "Name\tEndpoint\tSecure?\tAuthentication\n")
	bold := ansiterm.Styles(ansiterm.Bold)
	red := ansiterm.Foreground(ansiterm.Red)
	for _, ctx := range c.Contexts {
		if ctx.Name == c.CurrentContext {
			bold.Fprint(wt, ctx.Name)
		} else {
			fmt.Fprint(wt, ctx.Name)
		}
		_, _ = fmt.Fprintf(wt, "\t%s", ctx.Connection.Grpc)
		if ctx.Connection.Insecure {
			red.Fprint(wt, "\tplaintext")
		} else if ctx.Connection.Tls != nil && ctx.Connection.Tls.InsecureSkipVerify {
			red.Fprint(wt, "\tno verify")
		} else {
			fmt.Fprint(wt, "\tyes")
		}
		_, _ = fmt.Fprintf(wt, "\t%s\n", getAuthMethodString(ctx))
	}
	_ = wt.Flush()
}

func getAuthMethodString(ctx Context) string {
	if ctx.Identity.Token != "" {
		return "Token"
	}
	if ctx.Identity.TokenCommand != nil {
		return "Executable"
	}
	if ctx.Identity.Armory.ClientId != "" {
		return "Armory Cloud"
	}
	return "Unknown"
}
