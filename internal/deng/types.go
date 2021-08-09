package deng

import (
	"time"
)

type AppSummary struct {
	Name           string
	Deployments    int32
	LastSuccessful time.Time
	LastFailure    time.Time
}

type AccountSummary struct {
	Provider string
	Name     string
}

type DeploymentConfiguration struct {
	Application string
	EnvironmentType string
	EnvironmentName string
	EnvironmentNamespace string
	ViaAccount string
	ViaProvider string
	Version []string
	Kustomize bool
	Local bool
	Strategy string
	StrategySteps []string
	Wait bool
}