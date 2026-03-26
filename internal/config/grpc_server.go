package config

type GRPCServer struct {
	Address string `env:"GOPHKEEPER_GRPC_SERVER_ADDRESS" env-default:":50051"`
}
