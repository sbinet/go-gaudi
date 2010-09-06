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

// service
type Service struct {
	Component
}

func (svc *Service) Initialize() StatusCode {
	println(svc.CompName(), "initialize...")
	return StatusCode(0)
}

func (svc *Service) Finalize() StatusCode {
	println(svc.CompName(), "finalize...")
	return StatusCode(0)
}

func NewSvc(t,n string) IService {
	c := &Component{comp_name:n, comp_type:t}
	return &Service{*c}
}

// algtool
type AlgTool struct {
	Component
	parent IComponent
}

func (tool *AlgTool) CompName() string {
	return tool.parent.CompName() + "." + tool.Component.CompName()
}

func (tool *AlgTool) Initialize() StatusCode {
	println(tool.CompName(), "initialize...")
	return StatusCode(0)
}

func (tool *AlgTool) Finalize() StatusCode {
	println(tool.CompName(), "finalize...")
	return StatusCode(0)
}

func NewTool(t,n string, parent IComponent) IAlgTool {
	tool := &AlgTool{}
	tool.Component.comp_name = n
	tool.Component.comp_type = t
	tool.parent = parent
	return tool
}
/* EOF */
