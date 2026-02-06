package repository

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func TestDB(t *testing.T) (*sqlx.DB, func(...string)) {
	t.Helper()

	cfg, err := GetConfig()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode))
	if err != nil {
		t.Fatal()
	}

	if err = db.Ping(); err != nil {
		t.Fatal()
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS event (id SERIAL PRIMARY KEY, user_id INTEGER NOT NULL, description VARCHAR(255) NOT NULL,  date DATE NOT NULL, time TIME NOT NULL);")
	if err != nil {
		t.Fatal(err)
	}

	return db, func(tables ...string) {
		if len(tables) > 0 {
			_, err = db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
			if err != nil {
				t.Fatal(err)
			}
		}

		err = db.Close()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func GetConfig() (Config, error) {
	dir, err := getProjectRoot()
	dir = strings.Join([]string{dir, "18"}, "/")
	if err != nil {
		return Config{}, err
	}

	viper.AddConfigPath(fmt.Sprintf("%s\\configs", dir))
	viper.SetConfigName("config")
	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = godotenv.Load(fmt.Sprintf("%s\\.env", dir))
	if err != nil {
		return Config{}, err
	}

	return Config{
		Host:     viper.GetString("db.host"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASS"),
		DBName:   os.Getenv("POSTGRES_DB_TEST"),
		SSLMode:  viper.GetString("db.ssl_mode"),
	}, nil
}

func getProjectRoot() (string, error) {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error executing 'go list -m': %v, stderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}
