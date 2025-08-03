package config

type ServerConfig struct {
	AuthMethod  string
	KubeConfig  KubeConfig
	Destination string
}

type ExporterConfig struct {
}

type KubeConfig struct {
	Path string
}
