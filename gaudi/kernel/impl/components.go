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

type properties struct {
	props map[string]interface{}
}
func (p *properties) SetProperty(n string, v interface{}) StatusCode {
	p.props[n] = v
	return StatusCode(0)
}
func (p *properties) GetProperty(n string) interface{} {
	v,ok := p.props[n]
	if ok {
		return v
	}
	return nil
}
func (p *properties) GetProperties() []Property {
	props := make([]Property, len(p.props))
	i := 0
	for k,v := range p.props {
		props[i] = Property{Name:k, Value:v}
		i++
	}
	return props
}
// algorithm
type Algorithm struct {
	Component
	properties
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
	alg := &Algorithm{}
	alg.Component.comp_name = n
	alg.Component.comp_type = t
	return alg
}

// service
type Service struct {
	Component
	properties
}

func (svc *Service) InitializeSvc() StatusCode {
	println(svc.CompName(), "initialize...")
	return StatusCode(0)
}

func (svc *Service) FinalizeSvc() StatusCode {
	println(svc.CompName(), "finalize...")
	return StatusCode(0)
}

func NewSvc(t,n string) IService {
	svc := &Service{}
	svc.Component.comp_name = n
	svc.Component.comp_type = t
	return svc
}

// algtool
type AlgTool struct {
	Component
	properties
	parent IComponent
}

func (tool *AlgTool) CompName() string {
	return tool.parent.CompName() + "." + tool.Component.CompName()
}

func (tool *AlgTool) InitializeTool() StatusCode {
	println(tool.CompName(), "initialize...")
	return StatusCode(0)
}

func (tool *AlgTool) FinalizeTool() StatusCode {
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
