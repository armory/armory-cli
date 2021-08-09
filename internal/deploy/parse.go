package deploy

import (
	"fmt"
	"github.com/armory/armory-cli/internal/deng/protobuff"
	"github.com/spf13/pflag"
)

const (
	ParameterEnvironmentName      = "account"
	ParameterEnvironmentType      = "account-type"
	ParameterEnvironmentNamespace = "namespace"
	ParameterKustomize            = "kustomize"
	ParameterLocal                = "local"
	ParameterViaAccount           = "via-account"
	ParameterViaProvider          = "via-provider"
	ParameterApplication          = "app"
	ParameterWait                 = "wait"
	ParameterVersion              = "version"

	// Strategy flags
	ParameterStrategy      = "strategy"
	ParameterStrategySteps = "canary-step"
)

func NewParser(fs *pflag.FlagSet, args []string) *parser {
	return &parser{fs: fs, args: args, dep: &protobuff.Deployment{}}
}

type parser struct {
	fs       *pflag.FlagSet
	args     []string
	dep      *protobuff.Deployment
	versions map[string]string
}

func (p *parser) parse() (*protobuff.Deployment, error) {
	a, err := p.fs.GetString(ParameterApplication)
	if err != nil {
		return nil, err
	}
	p.dep.Application = a

	// Parse environment
	if err := p.parseEnvironment(); err != nil {
		return nil, err
	}

	switch p.dep.Environment.Provider {
	case protobuff.KubernetesProvider:
		// Parse artifacts now that we know the provider
		if err := p.parseKubernetesArtifacts(); err != nil {
			return nil, err
		}
	}

	if err := p.parseStrategy(); err != nil {
		return nil, err
	}

	return p.dep, nil
}

func (p *parser) parseEnvironment() error {
	t, err := p.fs.GetString(ParameterEnvironmentType)
	if err != nil {
		return err
	}
	switch t {
	case protobuff.KubernetesProvider:
		break
	default:
		return fmt.Errorf("unknown environment provider %s", t)
	}

	n, err := p.fs.GetString(ParameterEnvironmentName)
	if err != nil {
		return err
	}
	p.dep.Environment = &protobuff.Environment{
		Provider: t,
		Account:  n,
	}
	if t == protobuff.KubernetesProvider {
		ns, err := p.fs.GetString(ParameterEnvironmentNamespace)
		if err != nil {
			return err
		}
		p.dep.Environment.Qualifier = &protobuff.Environment_Kubernetes{
			Kubernetes: &protobuff.KubernetesQualifier{
				Namespace: ns,
			},
		}
	}
	return nil
}

func (p *parser) parseVia() (*protobuff.Via, error) {
	return nil, nil
}
