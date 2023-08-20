package moyai

import "github.com/df-mc/dragonfly/server"

type Config struct {
	server.UserConfig

	Moyai struct {
		Tebex string
	}

	Oomph struct {
		Port         int
		CombatMode   int
		MovementMode int
	}
}

func DefaultConfig() Config {
	c := Config{
		UserConfig: server.DefaultConfig(),
		Oomph: struct {
			Port         int
			CombatMode   int
			MovementMode int
		}{
			Port:         19132,
			CombatMode:   1,
			MovementMode: 1,
		},
	}
	return c
}
