package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/spf13/viper"
)

func Load(serviceName string, cfg interface{}) error {
	v := viper.New()

	v.SetConfigName(serviceName)
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/" + serviceName)

	v.SetEnvPrefix(strings.ToUpper(strings.ReplaceAll(serviceName, "-", "_")))
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("config: read file: %w", err)
		}
	}

	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("config: unmarshal: %w", err)
	}

	return nil
}

// LoadSecret fetches a secret from AWS Secrets Manager.
// Supports AWS_ENDPOINT_URL for LocalStack compatibility.
func LoadSecret(ctx context.Context, secretName string) (map[string]string, error) {
	opts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion("ap-south-1"),
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("config: load aws config: %w", err)
	}

	smOpts := func(o *secretsmanager.Options) {
		if endpoint := os.Getenv("AWS_ENDPOINT_URL"); endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	}

	client := secretsmanager.NewFromConfig(cfg, smOpts)
	result, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		return nil, fmt.Errorf("config: get secret %s: %w", secretName, err)
	}

	var secrets map[string]string
	if err := json.Unmarshal([]byte(aws.ToString(result.SecretString)), &secrets); err != nil {
		return nil, fmt.Errorf("config: parse secret %s: %w", secretName, err)
	}

	return secrets, nil
}
