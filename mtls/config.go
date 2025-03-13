package mtls

import "log/slog"

// Config contain the tls config passed by the config file.
type Config struct {
	// Hash is a unique hash of the cert + key + ca content.
	Hash string `json:"hash" mapstructure:"hash"`
	// Cert is the path to the TLS certificate.
	Cert string `json:"cert"     mapstructure:"cert"`
	// Key is the path to the TLS key.
	Key string `json:"key"      mapstructure:"key"`
	// CA is the path to the TLS CA certificate.
	Ca string `json:"ca"       mapstructure:"ca"`
	// Level TLS authentication level.
	Level Level `json:"level"    mapstructure:"level"`
	// Insecure is true if insecure TLS is allowed (client).
	Insecure bool `json:"insecure"    mapstructure:"insecure"`
}

func (cfg Config) SameAs(in Config) bool {
	return cfg.Hash == in.Hash &&
		cfg.Cert == in.Cert &&
		cfg.Key == in.Key &&
		cfg.Ca == in.Ca &&
		cfg.Level == in.Level
}

// // GetCert implemte Config.
// func (cfg Config) GetCert() string { return cfg.Cert }

// // GetKey implemte Config.
// func (cfg Config) GetKey() string { return cfg.Key }

// // GetKey implemte Config.
// func (cfg Config) GetCa() string { return cfg.Ca }

// func (cfg Config) GetLevel() Level { return cfg.Level }

// Empty implement Config.
func (cfg Config) Empty() bool {
	return cfg.Hash == "" && cfg.Cert == "" && cfg.Key == "" && !cfg.Insecure
}

func (cfg Config) AsAttrs() []any {
	if cfg.Empty() {
		return []any{}
	}

	return []any{
		slog.String("cert", cfg.Cert),
		slog.String("key", cfg.Key),
		slog.String("CA", cfg.Ca),
		slog.String("hash", cfg.Hash),

		slog.String("level", cfg.Level.String()),
		slog.Bool("insecure (client)", cfg.Insecure),
	}
}
