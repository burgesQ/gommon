package mtls

import "log/slog"

// Config contain the tls config passed by the config file.
type Config struct {
	Cert     string `json:"cert"     mapstructure:"cert"`
	Key      string `json:"key"      mapstructure:"key"`
	Ca       string `json:"ca"       mapstructure:"ca"`
	Level    Level  `json:"level"    mapstructure:"level"`
	Insecure bool   `json:"insecure"    mapstructure:"insecure"`
}

func (cfg Config) SameAs(in Config) bool {
	return cfg.Cert == in.Cert && cfg.Key == in.Key &&
		cfg.Ca == in.Ca && cfg.Level == in.Level
}

// // GetCert implemte Config.
// func (cfg Config) GetCert() string { return cfg.Cert }

// // GetKey implemte Config.
// func (cfg Config) GetKey() string { return cfg.Key }

// // GetKey implemte Config.
// func (cfg Config) GetCa() string { return cfg.Ca }

// func (cfg Config) GetLevel() Level { return cfg.Level }

// Empty implemte Config.
func (cfg Config) Empty() bool { return cfg.Cert == "" && cfg.Key == "" }

func (cfg Config) AsAttrs() []any {
	if cfg.Empty() {
		return []any{}
	}

	return []any{
		slog.String("cert", cfg.Cert), slog.String("key", cfg.Key),
		slog.String("CA", cfg.Ca), slog.String("level", cfg.Level.String()),
	}
}
