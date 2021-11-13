package server

type config struct {
	GRPCPort        int `value:"8080"`
	GRPCGatewayPort int `value:"8090"`
	DBConfig
}
