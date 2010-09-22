package kernel

import "fmt"
import "os"

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
	c := &Component{comp_name: n, comp_type: t}
	g_compsdb[n] = c
	return c
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

type OutputLevel int
const (
	LVL_VERBOSE OutputLevel = iota
	LVL_DEBUG
	LVL_INFO
	LVL_WARNING
	LVL_ERROR
	LVL_FATAL
	LVL_ALWAYS
	)

func (lvl OutputLevel) String() string {
	switch lvl {
	case LVL_VERBOSE: return "VERBOSE"
	case LVL_DEBUG:   return "DEBUG"
	case LVL_INFO:    return "INFO"
	case LVL_WARNING: return "WARNING"
	case LVL_ERROR:   return "ERROR"
	case LVL_FATAL:   return "FATAL"
	case LVL_ALWAYS:  return "ALWAYS"
	default:          return "???"
	}
	return "???"
}

// msgstream
type msgstream struct {
	name  string
	level OutputLevel
}


func (m *msgstream) SetOutputLevel(lvl OutputLevel) {
	m.level = lvl
}
func (m *msgstream) OutputLevel() OutputLevel {
	return m.level
}
func (m *msgstream) Msg(lvl OutputLevel, format string, a ...interface{}) (int, os.Error) {
	if m.level <= lvl {
		s := fmt.Sprintf(format, a)
		return fmt.Printf("%-10s %6s %s", m.name, lvl, s)
	}
	return 0, nil
}
func (m *msgstream) MsgVerbose(format string, a ...interface{}) (int, os.Error) {
	return m.Msg(LVL_VERBOSE, format, a)
}
func (m *msgstream) MsgDebug(format string, a ...interface{}) (int, os.Error) {
	return m.Msg(LVL_DEBUG, format, a)
}
func (m *msgstream) MsgInfo(format string, a ...interface{}) (int, os.Error) {
	return m.Msg(LVL_INFO, format, a)
}
func (m *msgstream) MsgWarning(format string, a ...interface{}) (int, os.Error) {
	return m.Msg(LVL_WARNING, format, a)
}
func (m *msgstream) MsgError(format string, a ...interface{}) (int, os.Error) {
	return m.Msg(LVL_ERROR, format, a)
}
func (m *msgstream) MsgFatal(format string, a ...interface{}) (int, os.Error) {
	return m.Msg(LVL_FATAL, format, a)
}
func (m *msgstream) MsgAlways(format string, a ...interface{}) (int, os.Error) {
	return m.Msg(LVL_ALWAYS, format, a)
}

// algorithm
type Algorithm struct {
	Component
	properties
	msgstream
}

func (alg *Algorithm) Initialize() StatusCode {
	alg.MsgInfo("initialize...\n")
	return StatusCode(0)
}

func (alg *Algorithm) Execute(ctx IEvtCtx) StatusCode {
	alg.MsgInfo("execute... [%v]\n", ctx)
	return StatusCode(0)
}

func (alg *Algorithm) Finalize() StatusCode {
	alg.MsgInfo("finalize...\n")
	return StatusCode(0)
}

func NewAlg(comp IComponent, t,n string) IAlgorithm {
	c := comp.(*Algorithm)
	c.Component.comp_name = n
	c.Component.comp_type = t
	c.properties.props = make(map[string]interface{})
	c.msgstream.name = n
	c.msgstream.level = LVL_INFO

	//g_compsdb[n] = c

	return c
}

// service
type Service struct {
	Component
	properties
	msgstream
}

func (svc *Service) InitializeSvc() StatusCode {
	svc.MsgInfo("initialize...\n")
	return StatusCode(0)
}

func (svc *Service) FinalizeSvc() StatusCode {
	svc.MsgInfo("finalize...\n")
	return StatusCode(0)
}

func NewSvc(comp IComponent, t,n string) IService {
	svc := comp.(*Service)
	svc.Component.comp_name = n
	svc.Component.comp_type = t
	svc.properties.props = make(map[string]interface{})

	svc.msgstream.name = n
	svc.msgstream.level = LVL_INFO

	//g_compsdb[n] = svc
	return svc
}

// algtool
type AlgTool struct {
	Component
	properties
	msgstream
	parent IComponent
}

func (tool *AlgTool) CompName() string {
	// FIXME: implement toolsvc !
	if tool.parent != nil {
		return tool.parent.CompName() + "." + tool.Component.CompName()
	}
	return "ToolSvc." + tool.Component.CompName()
}

func (tool *AlgTool) InitializeTool() StatusCode {
	tool.MsgInfo("initialize...\n")
	return StatusCode(0)
}

func (tool *AlgTool) FinalizeTool() StatusCode {
	tool.MsgInfo("finalize...\n")
	return StatusCode(0)
}

func NewTool(comp IComponent, t,n string, parent IComponent) IAlgTool {
	tool := comp.(*AlgTool)
	tool.Component.comp_name = n
	tool.Component.comp_type = t
	tool.properties.props = make(map[string]interface{})
	tool.msgstream = msgstream{name: tool.CompName(), level: LVL_INFO}
	tool.parent = parent

	//g_compsdb[n] = tool
	return tool
}

func init() {
	g_compsdb = make(comps_db)
	//fmt.Printf("--> components: %v\n", g_compsdb)
}

// checking implementations match interfaces
var _ = IAlgorithm(&Algorithm{})
var _ = IComponent(&Algorithm{})
var _ = IProperty(&Algorithm{})

var _ = IAlgTool(&AlgTool{})
var _ = IComponent(&AlgTool{})
var _ = IProperty(&AlgTool{})

var _ = IService(&Service{})
var _ = IComponent(&Service{})
var _ = IProperty(&Service{})

/* EOF */
