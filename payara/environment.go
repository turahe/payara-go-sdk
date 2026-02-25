package payara

import "github.com/turahe/payara-go-sdk/payara/types"

// Environment selects Sandbox or Production. Doc: Environmental Information
type Environment = types.Environment

const (
	EnvironmentSandbox    = types.EnvironmentSandbox
	EnvironmentProduction = types.EnvironmentProduction
)

// BaseURLForEnvironment returns the documented base URL for the environment.
// Sandbox: https://sandbox.payara.id:9090
// Production: https://openapi.payara.id:7654
func BaseURLForEnvironment(env Environment) string {
	switch env {
	case types.EnvironmentSandbox:
		return "https://sandbox.payara.id:9090"
	case types.EnvironmentProduction:
		return "https://openapi.payara.id:7654"
	default:
		return "https://openapi.payara.id:7654"
	}
}
