package micro

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

// Paths for packages used by code generated in this file,
// relative to the import_prefix of the generator.Generator.
const (
	contextPkgPath = "context"
	clientPkgPath  = "github.com/fananchong/v-micro/client"
	serverPkgPath  = "github.com/fananchong/v-micro/server"
)

func init() {
	generator.RegisterPlugin(new(micro))
}

// micro is an implementation of the Go protocol buffer compiler's
// plugin architecture.  It generates bindings for go-micro support.
type micro struct {
	gen *generator.Generator
}

// Name returns the name of this plugin, "vmicro".
func (g *micro) Name() string {
	return "vmicro"
}

// The names for packages imported in the generated code.
// They may vary from the final path component of the import path
// if the name is used by other packages.
var (
	contextPkg string
	clientPkg  string
	serverPkg  string
	pkgImports map[generator.GoPackageName]bool
)

// Init initializes the plugin.
func (g *micro) Init(gen *generator.Generator) {
	g.gen = gen
	contextPkg = generator.RegisterUniquePackageName("context", nil)
	clientPkg = generator.RegisterUniquePackageName("client", nil)
	serverPkg = generator.RegisterUniquePackageName("server", nil)
}

// Given a type name defined in a .proto, return its object.
// Also record that we're using it, to guarantee the associated import.
func (g *micro) objectNamed(name string) generator.Object {
	g.gen.RecordTypeUse(name)
	return g.gen.ObjectNamed(name)
}

// Given a type name defined in a .proto, return its name as we will print it.
func (g *micro) typeName(str string) string {
	return g.gen.TypeName(g.objectNamed(str))
}

// P forwards to g.gen.P.
func (g *micro) P(args ...interface{}) { g.gen.P(args...) }

// Generate generates code for the services in the given file.
func (g *micro) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	for i, service := range file.FileDescriptorProto.Service {
		g.generateService(file, service, i)
	}
}

// GenerateImports generates the import declaration for this file.
func (g *micro) GenerateImports(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	g.P("import (")
	g.P(contextPkg, " ", strconv.Quote(path.Join(g.gen.ImportPrefix, contextPkgPath)))
	g.P(clientPkg, " ", strconv.Quote(path.Join(g.gen.ImportPrefix, clientPkgPath)))
	g.P(serverPkg, " ", strconv.Quote(path.Join(g.gen.ImportPrefix, serverPkgPath)))
	g.P(")")
	g.P()
}

// reservedClientName records whether a client name is reserved on the client side.
var reservedClientName = map[string]bool{
	// TODO: do we need any in go-micro?
}

func unexport(s string) string {
	if len(s) == 0 {
		return ""
	}
	name := strings.ToLower(s[:1]) + s[1:]
	if pkgImports[generator.GoPackageName(name)] {
		return name + "_"
	}
	return name
}

// generateService generates all the code for the named service.
func (g *micro) generateService(file *generator.FileDescriptor, service *pb.ServiceDescriptorProto, index int) {
	path := fmt.Sprintf("6,%d", index) // 6 means service.

	origServName := service.GetName()
	serviceName := strings.ToLower(service.GetName())
	if pkg := file.GetPackage(); pkg != "" {
		serviceName = pkg
	}
	servName := generator.CamelCase(origServName)
	servAlias := servName + "Service"

	// strip suffix
	if strings.HasSuffix(servAlias, "ServiceService") {
		servAlias = strings.TrimSuffix(servAlias, "Service")
	}

	g.P()
	g.P("// Client API for ", servName, " service")
	g.P()

	// Client interface.
	g.P("type ", servAlias, " interface {")
	for i, method := range service.Method {
		g.gen.PrintComments(fmt.Sprintf("%s,2,%d", path, i)) // 2 means method in a service.
		g.P(g.generateClientSignature(servName, method))
	}
	g.P("}")
	g.P()

	// Client structure.
	g.P("type ", unexport(servAlias), " struct {")
	g.P("c ", clientPkg, ".Client")
	g.P("name string")
	g.P("}")
	g.P()

	// NewClient factory.
	g.P("func New", servAlias, " (name string, c ", clientPkg, ".Client) ", servAlias, " {")
	g.P("if c == nil {")
	g.P("panic(\"client is nil\")")
	g.P("}")
	g.P("if len(name) == 0 {")
	g.P("panic(\"name is nil\")")
	g.P("}")
	g.P("return &", unexport(servAlias), "{")
	g.P("c: c,")
	g.P("name: name,")
	g.P("}")
	g.P("}")
	g.P()
	var methodIndex, streamIndex int
	serviceDescVar := "_" + servName + "_serviceDesc"
	// Client method implementations.
	for _, method := range service.Method {
		var descExpr string
		if !method.GetServerStreaming() {
			// Unary RPC method
			descExpr = fmt.Sprintf("&%s.Methods[%d]", serviceDescVar, methodIndex)
			methodIndex++
		} else {
			// Streaming RPC method
			descExpr = fmt.Sprintf("&%s.Streams[%d]", serviceDescVar, streamIndex)
			streamIndex++
		}
		g.generateClientMethod(serviceName, servName, serviceDescVar, method, descExpr)
	}

	g.P("// Server API for ", servName, " service")
	g.P()

	// Server interface.
	serverType := servName + "Handler"
	g.P("type ", serverType, " interface {")
	for i, method := range service.Method {
		g.gen.PrintComments(fmt.Sprintf("%s,2,%d", path, i)) // 2 means method in a service.
		g.P(g.generateServerSignature(servName, method))
	}
	g.P("}")
	g.P()

	// Server registration.
	g.P("func Register", servName, "Handler(s ", serverPkg, ".Server, hdlr ", serverType, ") error {")
	g.P("return s.Handle(hdlr)")
	g.P("}")
	g.P()
}

// generateClientSignature returns the client-side signature for a method.
func (g *micro) generateClientSignature(servName string, method *pb.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)
	if reservedClientName[methName] {
		methName += "_"
	}
	reqArg := ", in *" + g.typeName(method.GetInputType())
	if method.GetClientStreaming() {
		reqArg = ""
	}
	return fmt.Sprintf("%s(ctx %s.Context%s, opts ...%s.CallOption) error", methName, contextPkg, reqArg, clientPkg)
}

func (g *micro) generateClientMethod(reqServ, servName, serviceDescVar string, method *pb.MethodDescriptorProto, descExpr string) {
	reqMethod := fmt.Sprintf("%s.%s", servName, method.GetName())
	methName := generator.CamelCase(method.GetName())
	inType := g.typeName(method.GetInputType())
	outType := g.typeName(method.GetOutputType())

	servAlias := servName + "Service"

	// strip suffix
	if strings.HasSuffix(servAlias, "ServiceService") {
		servAlias = strings.TrimSuffix(servAlias, "Service")
	}

	g.P("func (c *", unexport(servAlias), ") ", g.generateClientSignature(servName, method), "{")
	if !method.GetServerStreaming() && !method.GetClientStreaming() {
		g.P(`req := c.c.NewRequest(c.name, "`, reqMethod, `", in)`)
		// TODO: Pass descExpr to Invoke.
		g.P("err := ", `c.c.Call(ctx, req, opts...)`)
		g.P("if err != nil { return err }")
		g.P("return nil")
		g.P("}")
		g.P()
		return
	}
	streamType := unexport(servAlias) + methName
	g.P(`req := c.c.NewRequest(c.name, "`, reqMethod, `", &`, inType, `{})`)
	g.P("stream, err := c.c.Stream(ctx, req, opts...)")
	g.P("if err != nil { return nil, err }")

	if !method.GetClientStreaming() {
		g.P("if err := stream.Send(in); err != nil { return nil, err }")
	}

	g.P("return &", streamType, "{stream}, nil")
	g.P("}")
	g.P()

	genSend := method.GetClientStreaming()
	genRecv := method.GetServerStreaming()

	// Stream auxiliary types and methods.
	g.P("type ", servName, "_", methName, "Service interface {")
	g.P("SendMsg(interface{}) error")
	g.P("RecvMsg(interface{}) error")
	g.P("Close() error")

	if genSend {
		g.P("Send(*", inType, ") error")
	}
	if genRecv {
		g.P("Recv() (*", outType, ", error)")
	}
	g.P("}")
	g.P()

	g.P("type ", streamType, " struct {")
	g.P("stream ", clientPkg, ".Stream")
	g.P("}")
	g.P()

	g.P("func (x *", streamType, ") Close() error {")
	g.P("return x.stream.Close()")
	g.P("}")
	g.P()

	g.P("func (x *", streamType, ") SendMsg(m interface{}) error {")
	g.P("return x.stream.Send(m)")
	g.P("}")
	g.P()

	g.P("func (x *", streamType, ") RecvMsg(m interface{}) error {")
	g.P("return x.stream.Recv(m)")
	g.P("}")
	g.P()

	if genSend {
		g.P("func (x *", streamType, ") Send(m *", inType, ") error {")
		g.P("return x.stream.Send(m)")
		g.P("}")
		g.P()

	}

	if genRecv {
		g.P("func (x *", streamType, ") Recv() (*", outType, ", error) {")
		g.P("m := new(", outType, ")")
		g.P("err := x.stream.Recv(m)")
		g.P("if err != nil {")
		g.P("return nil, err")
		g.P("}")
		g.P("return m, nil")
		g.P("}")
		g.P()
	}
}

// generateServerSignature returns the server-side signature for a method.
func (g *micro) generateServerSignature(servName string, method *pb.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)
	if reservedClientName[methName] {
		methName += "_"
	}

	var reqArgs []string
	ret := "error"
	reqArgs = append(reqArgs, contextPkg+".Context")

	if !method.GetClientStreaming() {
		reqArgs = append(reqArgs, "*"+g.typeName(method.GetInputType()))
	}
	if method.GetServerStreaming() || method.GetClientStreaming() {
		reqArgs = append(reqArgs, servName+"_"+generator.CamelCase(origMethName)+"Stream")
	}
	if !method.GetClientStreaming() && !method.GetServerStreaming() {
		reqArgs = append(reqArgs, "*"+g.typeName(method.GetOutputType()))
	}
	return methName + "(" + strings.Join(reqArgs, ", ") + ") " + ret
}

// AddPluginToParams Simplify the protoc call statement by adding 'plugins=vmicro' directly to the command line arguments.
func AddPluginToParams(p string) string {
	params := p
	if strings.Contains(params, "plugins=") {
		params = strings.Replace(params, "plugins=", "plugins=vmicro+", -1)
	} else {
		if len(params) > 0 {
			params += ","
		}
		params += "plugins=vmicro"
	}
	return params
}
