package kernel

type Component struct {
	comp_name string
	comp_type string
}

func (c *Component) CompName() string {
	return c.comp_name
}

func (c *Component) CompType() string {
	return c.comp_type
}

func NewComponent(t,n string) IComponent {
	return &Component{comp_name: n, comp_type: t}
}

// algorithm
type Algorithm struct {
	Component
}

func (alg *Algorithm) Initialize() StatusCode {
	println(alg.CompName(), "initialize...")
	return StatusCode(0)
}

func (alg *Algorithm) Execute() StatusCode {
	println(alg.CompName(), "execute...")
	return StatusCode(0)
}

func (alg *Algorithm) Finalize() StatusCode {
	println(alg.CompName(), "finalize...")
	return StatusCode(0)
}

func NewAlg(t,n string) IAlgorithm {
	c := &Component{comp_name:n, comp_type:t}
	return &Algorithm{*c}
}
/* EOF */
