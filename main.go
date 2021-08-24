package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

func main() {
	if !clowder.IsClowderEnabled() {
		fmt.Fprintln(os.Stderr, "Clowder not enabled - exiting")
		os.Exit(1)
	}

	cfg := clowder.LoadedConfig
	cmd := exec.Command(os.Args[1], os.Args[2:]...)

	// run sub-command with same stdin/stdout
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// building out a big ole' slice of ENV vars, starting with the current env.
	env := os.Environ()
	env = append(env, dbEnv(cfg)...)
	env = append(env, inMemoryEnv(cfg)...)
	env = append(env, clowdwatchEnv(cfg)...)
	env = append(env, kafkaEnv(cfg))
	env = append(env, kafkaTopics(cfg)...)

	cmd.Env = env

	// pass go, collect $200
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command: %v\n", err)
		os.Exit(1)
	}
}

func dbEnv(cfg *clowder.AppConfig) []string {
	return []string{
		fmt.Sprintf("DATABASE_HOST=%s", cfg.Database.Hostname),
		fmt.Sprintf("DATABASE_PORT=%d", cfg.Database.Port),
		fmt.Sprintf("PGSSLMODE=%s", cfg.Database.SslMode),
		fmt.Sprintf("DATABASE_NAME=%s", cfg.Database.Name),
		fmt.Sprintf("DATABASE_USER=%s", cfg.Database.Username),
		fmt.Sprintf("DATABASE_PASSWORD=%s", cfg.Database.Password),
		fmt.Sprintf("DATABASE_ADMIN_USER=%s", cfg.Database.AdminUsername),
		fmt.Sprintf("DATABASE_ADMIN_PASSWORD=%s", cfg.Database.AdminPassword),
	}
}

func inMemoryEnv(cfg *clowder.AppConfig) []string {
	// these are always there
	env := []string{
		fmt.Sprintf("IN_MEMORY_HOST=%s", cfg.InMemoryDb.Hostname),
		fmt.Sprintf("IN_MEMORY_PORT=%d", cfg.InMemoryDb.Port),
	}

	// can be nil
	if cfg.InMemoryDb.Username != nil {
		env = append(env, fmt.Sprintf("IN_MEMORY_USER=%s", *cfg.InMemoryDb.Username))
	}

	// can be nil
	if cfg.InMemoryDb.Password != nil {
		env = append(env, fmt.Sprintf("IN_MEMORY_PASSWORD=%s", *cfg.InMemoryDb.Password))
	}

	return env
}

func kafkaEnv(cfg *clowder.AppConfig) string {
	brokers := make([]string, len(cfg.Kafka.Brokers))
	for i, broker := range cfg.Kafka.Brokers {
		brokers[i] = fmt.Sprintf("%s:%d", broker.Hostname, *broker.Port)
	}

	return fmt.Sprintf("KAFKA_BROKERS=%s", strings.Join(brokers, ","))
}

func kafkaTopics(cfg *clowder.AppConfig) []string {
	topics := make([]string, len(cfg.Kafka.Topics))

	for i, t := range cfg.Kafka.Topics {
		topics[i] = fmt.Sprintf("TOPIC_NAME_%v=%v", t.RequestedName, t.Name)
	}

	return topics
}

func clowdwatchEnv(cfg *clowder.AppConfig) []string {
	return []string{
		fmt.Sprintf("CW_ACCESS_KEY_ID=%s", cfg.Logging.Cloudwatch.AccessKeyId),
		fmt.Sprintf("CW_SECRET_ACCESS_KEY=%s", cfg.Logging.Cloudwatch.SecretAccessKey),
		fmt.Sprintf("CW_LOG_GROUP=%s", cfg.Logging.Cloudwatch.LogGroup),
		fmt.Sprintf("CW_REGION=%s", cfg.Logging.Cloudwatch.Region),
	}
}
