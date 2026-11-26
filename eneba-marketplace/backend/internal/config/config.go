package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct{ DatabaseURL, Port, FrontendURL string }

func Load() Config {
	if f, err := os.Open(".env"); err == nil {
		defer f.Close()
		s := bufio.NewScanner(f)
		for s.Scan() {
			k, v, ok := strings.Cut(s.Text(), "=")
			k = strings.TrimSpace(k)
			if ok && k != "" && !strings.HasPrefix(k, "#") {
				if _, exists := os.LookupEnv(k); !exists {
					os.Setenv(k, strings.TrimSpace(v))
				}
			}
		}
	}
	get := func(k, d string) string {
		if v := os.Getenv(k); v != "" {
			return v
		}
		return d
	}
	return Config{os.Getenv("DATABASE_URL"), get("PORT", "8080"), get("FRONTEND_URL", "http://localhost:5173")}
}
