package environment

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

func (e Environment) IsProduction() bool {
	return e == Production
}
