package config

type Config struct {
	Database struct {
		Host     string `env:"DATABASE_HOST" env-default:"localhost"`
		Port     int    `env:"DATABASE_PORT" env-default:"5432"`
		Name     string `env:"DATABASE_NAME" env-default:"book_manager"`
		Username string `env:"DATABASE_USERNAME" env-default:"admin"`
		Password string `env:"DATABASE_PASSWORD" env-default:"admin"`
	}
}
