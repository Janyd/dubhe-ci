package config

type (
	Config struct {
		Workspace string     `toml:"workspace"`
		Capacity  int        `toml:"capacity"`
		Http      HttpServer `toml:"http"`
		Server    GrpcServer `toml:"server"`
		Database  Database   `toml:"database"`
		Log       Log        `toml:"log"`
		JWTAuth   JWTAuth    `toml:"jwt"`
		SuperUser SuperUser  `toml:"super_user"`
	}

	HttpServer struct {
		Address string
	}

	GrpcServer struct {
		Address string `toml:"address"`
		Network string `toml:"network"`
	}

	Database struct {
		Debug        bool   `toml:"debug"`
		DBType       string `toml:"db_type"`
		DSN          string `toml:"dsn"`
		MaxLifetime  int    `toml:"max_lifetime"`
		MaxOpenConns int    `toml:"max_open_conns"`
		MaxIdleConns int    `toml:"max_idle_conns"`
	}

	Log struct {
		Level     int    `toml:"level"`
		Format    string `toml:"format"`
		InfoPath  string `toml:"info_path"`
		ErrorPath string `toml:"error_path"`
	}

	JWTAuth struct {
		SigningKey string `toml:"signing_key"`
		Expired    int    `toml:"expired"`
	}

	SuperUser struct {
		Nickname string `toml:"nickname"`
		Username string `toml:"username"`
		Password string `toml:"password"`
	}
)
