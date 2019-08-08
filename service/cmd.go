package service

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fananchong/v-micro/client"
	crpc "github.com/fananchong/v-micro/client/rpc"
	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/common/log/logrus"
	"github.com/fananchong/v-micro/registry"
	"github.com/fananchong/v-micro/registry/mdns"
	"github.com/fananchong/v-micro/selector"
	"github.com/fananchong/v-micro/selector/cache"
	"github.com/fananchong/v-micro/server"
	srpc "github.com/fananchong/v-micro/server/rpc"
	"github.com/fananchong/v-micro/transport"
	"github.com/fananchong/v-micro/transport/gotcp"
	"github.com/micro/cli"
)

// Cmd cmd
type Cmd interface {
	// The cli app within this cmd
	App() *cli.App
	// Adds options, parses flags and initialise
	// exits on error
	Init(opts ...Option) error
	// Options set within this command
	Options() Options
}

type cmd struct {
	opts Options
	app  *cli.App
}

var (
	// DefaultCmd default cmd
	DefaultCmd = newCmd()

	// DefaultFlags default flags
	DefaultFlags = []cli.Flag{
		cli.StringFlag{
			Name:   "logger",
			EnvVar: "MICRO_LOGGER",
			Usage:  "Logger for v-micro; logrus",
		},
		cli.StringFlag{
			Name:   "client",
			EnvVar: "MICRO_CLIENT",
			Usage:  "Client for v-micro; rpc",
		},
		cli.IntFlag{
			Name:   "client_retries",
			EnvVar: "MICRO_CLIENT_RETRIES",
			Value:  client.DefaultRetries,
			Usage:  "Sets the client retries. Default: 1",
		},
		cli.IntFlag{
			Name:   "register_ttl",
			EnvVar: "MICRO_REGISTER_TTL",
			Usage:  "Register TTL in seconds",
		},
		cli.IntFlag{
			Name:   "register_interval",
			EnvVar: "MICRO_REGISTER_INTERVAL",
			Usage:  "Register interval in seconds",
		},
		cli.StringFlag{
			Name:   "server",
			EnvVar: "MICRO_SERVER",
			Usage:  "Server for v-micro; rpc",
		},
		cli.StringFlag{
			Name:   "server_name",
			EnvVar: "MICRO_SERVER_NAME",
			Usage:  "Name of the server",
		},
		cli.StringFlag{
			Name:   "server_version",
			EnvVar: "MICRO_SERVER_VERSION",
			Usage:  "Version of the server. 1.1.0",
		},
		cli.StringFlag{
			Name:   "server_id",
			EnvVar: "MICRO_SERVER_ID",
			Usage:  "Id of the server. Auto-generated if not specified",
		},
		cli.StringFlag{
			Name:   "server_address",
			EnvVar: "MICRO_SERVER_ADDRESS",
			Usage:  "Bind address for the server. 127.0.0.1:8080",
		},
		cli.StringFlag{
			Name:   "server_advertise",
			EnvVar: "MICRO_SERVER_ADVERTISE",
			Usage:  "Used instead of the server_address when registering with discovery. 127.0.0.1:8080",
		},
		cli.StringSliceFlag{
			Name:   "server_metadata",
			EnvVar: "MICRO_SERVER_METADATA",
			Value:  &cli.StringSlice{},
			Usage:  "A list of key-value pairs defining metadata. version=1.0.0",
		},
		cli.StringFlag{
			Name:   "registry",
			EnvVar: "MICRO_REGISTRY",
			Usage:  "Registry for discovery. consul, mdns",
		},
		cli.StringFlag{
			Name:   "registry_address",
			EnvVar: "MICRO_REGISTRY_ADDRESS",
			Usage:  "Comma-separated list of registry addresses",
		},
		cli.StringFlag{
			Name:   "selector",
			EnvVar: "MICRO_SELECTOR",
			Usage:  "Selector used to pick nodes for querying",
		},
		cli.StringFlag{
			Name:   "transport",
			EnvVar: "MICRO_TRANSPORT",
			Usage:  "Transport mechanism used; http",
		},
		cli.IntFlag{
			Name:   "log_level",
			EnvVar: "MICRO_LOG_LEVEL",
			Usage:  "Logger level",
		},
	}

	// DefaultLogs default logs
	DefaultLogs = map[string]func(...log.Option) log.Logger{
		"logrus": logrus.NewLogger,
	}

	// DefaultClients default clients
	DefaultClients = map[string]func(...client.Option) client.Client{
		"rpc": crpc.NewClient,
	}

	// DefaultRegistries default registries
	DefaultRegistries = map[string]func(...registry.Option) registry.Registry{
		"mdns": mdns.NewRegistry,
	}

	// DefaultSelectors default selectors
	DefaultSelectors = map[string]func(...selector.Option) selector.Selector{
		"cache": cache.NewSelector,
	}

	// DefaultServers default servers
	DefaultServers = map[string]func(...server.Option) server.Server{
		"rpc": srpc.NewServer,
	}

	// DefaultTransports default transports
	DefaultTransports = map[string]func(...transport.Option) transport.Transport{
		"gotcp": gotcp.NewTransport,
	}

	// used for default selection as the fall back
	defaultLog       = "logrus"
	defaultClient    = "rpc"
	defaultServer    = "rpc"
	defaultRegistry  = "mdns"
	defaultSelector  = "cache"
	defaultTransport = "gotcp"
)

func init() {
	rand.Seed(time.Now().Unix())
	help := cli.HelpPrinter
	cli.HelpPrinter = func(writer io.Writer, templ string, data interface{}) {
		help(writer, templ, data)
		os.Exit(0)
	}
}

func newCmd(opts ...Option) Cmd {
	if log.DefaultLogger == nil {
		log.DefaultLogger = DefaultLogs[defaultLog]()
	}
	if registry.DefaultRegistry == nil {
		registry.DefaultRegistry = DefaultRegistries[defaultRegistry]()
	}
	if selector.DefaultSelector == nil {
		selector.DefaultSelector = DefaultSelectors[defaultSelector]()
	}
	if transport.DefaultTransport == nil {
		transport.DefaultTransport = DefaultTransports[defaultTransport]()
	}
	if client.DefaultClient == nil {
		client.DefaultClient = DefaultClients[defaultClient]()
	}
	if server.DefaultServer == nil {
		server.DefaultServer = DefaultServers[defaultServer]()
	}
	options := Options{
		Logger:    &log.DefaultLogger,
		Client:    &client.DefaultClient,
		Registry:  &registry.DefaultRegistry,
		Server:    &server.DefaultServer,
		Selector:  &selector.DefaultSelector,
		Transport: &transport.DefaultTransport,

		Loggers:    DefaultLogs,
		Clients:    DefaultClients,
		Registries: DefaultRegistries,
		Selectors:  DefaultSelectors,
		Servers:    DefaultServers,
		Transports: DefaultTransports,
	}

	for _, o := range opts {
		o(&options)
	}

	if len(options.Description) == 0 {
		options.Description = "a v-micro service"
	}

	cmd := new(cmd)
	cmd.opts = options
	cmd.app = cli.NewApp()
	cmd.app.Name = cmd.opts.Name
	cmd.app.Version = cmd.opts.Version
	cmd.app.Usage = cmd.opts.Description
	cmd.app.Before = cmd.Before
	cmd.app.Flags = DefaultFlags
	cmd.app.Action = func(c *cli.Context) {}

	if len(options.Version) == 0 {
		cmd.app.HideVersion = true
	}

	return cmd
}

func (c *cmd) App() *cli.App {
	return c.app
}

func (c *cmd) Options() Options {
	return c.opts
}

func (c *cmd) Before(ctx *cli.Context) error {
	// If flags are set then use them otherwise do nothing
	var serverOpts []server.Option
	var clientOpts []client.Option
	var logOpts []log.Option

	// Set the logger
	if name := ctx.String("logger"); len(name) > 0 {
		// only change if we have the client and type differs
		if l, ok := c.opts.Loggers[name]; ok && (*c.opts.Logger).String() != name {
			*c.opts.Logger = l()
		}
	}

	// Set the client
	if name := ctx.String("client"); len(name) > 0 {
		// only change if we have the client and type differs
		if cl, ok := c.opts.Clients[name]; ok && (*c.opts.Client).String() != name {
			*c.opts.Client = cl()
		}
	}

	// Set the server
	if name := ctx.String("server"); len(name) > 0 {
		// only change if we have the server and type differs
		if s, ok := c.opts.Servers[name]; ok && (*c.opts.Server).String() != name {
			*c.opts.Server = s()
		}
	}

	// Set the registry
	if name := ctx.String("registry"); len(name) > 0 && (*c.opts.Registry).String() != name {
		r, ok := c.opts.Registries[name]
		if !ok {
			return fmt.Errorf("Registry %s not found", name)
		}

		*c.opts.Registry = r()
		serverOpts = append(serverOpts, server.Registry(*c.opts.Registry))
		clientOpts = append(clientOpts, client.Registry(*c.opts.Registry))

		if err := (*c.opts.Selector).Init(selector.Registry(*c.opts.Registry)); err != nil {
			log.Errorf("Error configuring registry: %v", err)
			os.Exit(1)
		}

		clientOpts = append(clientOpts, client.Selector(*c.opts.Selector))
	}

	// Set the selector
	if name := ctx.String("selector"); len(name) > 0 && (*c.opts.Selector).String() != name {
		s, ok := c.opts.Selectors[name]
		if !ok {
			return fmt.Errorf("Selector %s not found", name)
		}

		*c.opts.Selector = s(selector.Registry(*c.opts.Registry))

		// No server option here. Should there be?
		clientOpts = append(clientOpts, client.Selector(*c.opts.Selector))
	}

	// Set the transport
	if name := ctx.String("transport"); len(name) > 0 && (*c.opts.Transport).String() != name {
		t, ok := c.opts.Transports[name]
		if !ok {
			return fmt.Errorf("Transport %s not found", name)
		}

		*c.opts.Transport = t()
		serverOpts = append(serverOpts, server.Transport(*c.opts.Transport))
		clientOpts = append(clientOpts, client.Transport(*c.opts.Transport))
	}

	// Parse the server options
	metadata := make(map[string]string)
	for _, d := range ctx.StringSlice("server_metadata") {
		var key, val string
		parts := strings.Split(d, "=")
		key = parts[0]
		if len(parts) > 1 {
			val = strings.Join(parts[1:], "=")
		}
		metadata[key] = val
	}

	if len(metadata) > 0 {
		serverOpts = append(serverOpts, server.Metadata(metadata))
	}

	if len(ctx.String("registry_address")) > 0 {
		if err := (*c.opts.Registry).Init(registry.Addrs(strings.Split(ctx.String("registry_address"), ",")...)); err != nil {
			log.Errorf("Error configuring registry: %v", err)
			os.Exit(1)
		}
	}
	serverName := server.DefaultName
	if len(ctx.String("server_name")) > 0 {
		serverOpts = append(serverOpts, server.Name(ctx.String("server_name")))
		serverName = ctx.String("server_name")
	}

	if len(ctx.String("server_version")) > 0 {
		serverOpts = append(serverOpts, server.Version(ctx.String("server_version")))
	}

	serverID := server.DefaultID
	if len(ctx.String("server_id")) > 0 {
		serverOpts = append(serverOpts, server.ID(ctx.String("server_id")))
		serverID = ctx.String("server_id")
	}

	if len(ctx.String("server_address")) > 0 {
		serverOpts = append(serverOpts, server.Address(ctx.String("server_address")))
	}

	if len(ctx.String("server_advertise")) > 0 {
		serverOpts = append(serverOpts, server.Advertise(ctx.String("server_advertise")))
	}

	if ttl := time.Duration(ctx.GlobalInt("register_ttl")); ttl > 0 {
		serverOpts = append(serverOpts, server.RegisterTTL(ttl*time.Second))
	}

	if val := time.Duration(ctx.GlobalInt("register_interval")); val > 0 {
		serverOpts = append(serverOpts, server.RegisterInterval(val*time.Second))
	}

	// client opts
	if r := ctx.Int("client_retries"); r >= 0 {
		clientOpts = append(clientOpts, client.Retries(r))
	}

	if level := ctx.Int("log_level"); level >= 0 {
		logOpts = append(logOpts, log.Level(log.LevelType(level)))
	}
	logOpts = append(logOpts, log.Name(fmt.Sprintf("%s_%s", serverName, serverID)))

	if err := (*c.opts.Logger).Init(logOpts...); err != nil {
		log.Errorf("Error init logger: %v", err)
		os.Exit(1)
	}
	log.Infof("Logger [%s] enable", (*c.opts.Logger).String())

	// We have some command line opts for the server.
	// Lets set it up
	if len(serverOpts) > 0 {
		if err := (*c.opts.Server).Init(serverOpts...); err != nil {
			log.Errorf("Error configuring server: %v", err)
			os.Exit(1)
		}
	}

	// Use an init option?
	if len(clientOpts) > 0 {
		if err := (*c.opts.Client).Init(clientOpts...); err != nil {
			log.Errorf("Error configuring client: %v", err)
			os.Exit(1)
		}
	}

	return nil
}

func (c *cmd) Init(opts ...Option) error {
	for _, o := range opts {
		o(&c.opts)
	}
	c.app.Name = c.opts.Name
	c.app.Version = c.opts.Version
	c.app.HideVersion = len(c.opts.Version) == 0
	c.app.Usage = c.opts.Description
	c.app.Flags = append(c.app.Flags, c.opts.Flags...)
	c.app.Action = c.opts.Action
	c.app.RunAndExitOnError()
	return nil
}

// DefaultOptions default options
func DefaultOptions() Options {
	return DefaultCmd.Options()
}

// App app
func App() *cli.App {
	return DefaultCmd.App()
}

// Init init
func Init(opts ...Option) error {
	return DefaultCmd.Init(opts...)
}

// NewCmd new
func NewCmd(opts ...Option) Cmd {
	return newCmd(opts...)
}
