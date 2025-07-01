package config

import (
	"errors"
	string_helper "financing-offer/pkg/string-helper"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/fs"
	"github.com/knadh/koanf/v2"
	iofs "io/fs"
	"strings"
)

const (
	EnvPrefix     = "APP__"
	EnvProduction = "prod"

	SavingsTaskQueueName = "savings_queue"
)

type AppConfig struct {
	Env         string   `koanf:"env"`
	HttpPort    int      `koanf:"httpPort"`
	ConnectPort int      `koanf:"connectPort"`
	Db          DbConfig `koanf:"db"`
	Jwt         struct {
		PublicKey string `koanf:"publicKey"`
	} `koanf:"jwt"`
	ModelGeneration struct {
		Path             string   `koanf:"path"`
		IgnoredTables    []string `koanf:"ignoredTables"`
		ImmutableColumns []string `koanf:"immutableColumns"`
	} `koanf:"modelGeneration"`
	Kafka      KafkaConfig `koanf:"kafka"`
	Mattermost struct {
		WebhookUrl string `koanf:"webhookUrl"`
	} `koanf:"mattermost"`
	Cron              Cron                     `koanf:"cron"`
	Temporal          TemporalClientConfig     `koanf:"temporal"`
	FinancialProduct  FinancialProductConfig   `koanf:"financialProduct"`
	MoService         MoServiceConfig          `koanf:"moService"`
	Features          map[string]FeatureConfig `koanf:"features"`
	LoanRequest       LoanRequestConfig        `koanf:"loanRequest"`
	FinancingApi      FinancingApiConfig       `koanf:"financingApi"`
	BestPromotions    BestPromotionsConfig     `koanf:"bestPromotions"`
	OrderService      OrderServiceConfig       `koanf:"orderService"`
	FlexOpenApi       FlexOpenApiConfig        `koanf:"flexOpenApi"`
	OdooService       OdooServiceConfig        `koanf:"OdooService"`
	ProductCategoryId int64                    `koanf:"productCategoryId"`
	OdooCategoryId    int64                    `koanf:"odooCategoryId"`
}

type LoanRequestConfig struct {
	ExpireDays                   int     `koanf:"expireDays"`
	MaxGuaranteedDuration        int     `koanf:"maxGuaranteedDuration"`
	GuaranteeFeeRate             float64 `koanf:"guaranteeFeeRate"`
	MinimumAppVersion            string  `koanf:"minimumAppVersion"`
	MinimumAppVersionDerivative  string  `koanf:"minimumAppVersionDerivative"`
	DeclinedRequestDisplayPeriod int     `koanf:"declinedRequestDisplayPeriod"`
}

type FeatureConfig struct {
	Enable      bool     `koanf:"enable"`
	InvestorIds []string `koanf:"investorIds"`
}

type FinancialProductConfig struct {
	Url   string `koanf:"url"`
	Token string `koanf:"token"`
}

type MoServiceConfig struct {
	Url   string `koanf:"url"`
	Token string `koanf:"token"`
}

type FinancingApiConfig struct {
	Url   string `koanf:"url"`
	Token string `koanf:"token"`
}

type OrderServiceConfig struct {
	Url   string `koanf:"url"`
	Token string `koanf:"token"`
}

type OdooServiceConfig struct {
	Url      string `koanf:"url"`
	Db       string `koanf:"db"`
	Uid      string `koanf:"uid"`
	Password string `koanf:"password"`
}

type BestPromotionsConfig struct {
	LoanPackageIds []int64 `koanf:"loanPackageIds"`
}

type DbConfig struct {
	User        string `koanf:"user"`
	Password    string `koanf:"password"`
	DbName      string `koanf:"dbName"`
	Port        string `koanf:"port"`
	Host        string `koanf:"host"`
	EnableSsl   bool   `koanf:"enableSsl"`
	AutoMigrate bool   `koanf:"autoMigrate"`
}

type KafkaConfig struct {
	Host              string `koanf:"host"`
	Retry             int    `koanf:"retry"`
	NotificationTopic string `koanf:"notificationTopic"`
}

type TemporalClientConfig struct {
	Host      string `koanf:"host"`
	Namespace string `koanf:"namespace"`
}

type Cron struct {
	ExpireLoanOffers    string `koanf:"expireLoanOffers"`
	DeclineLoanRequests string `koanf:"declineLoanRequests"`
}

type MarginPoolConfig struct {
	Ids []int64 `koanf:"ids"`
}

type FlexOpenApiConfig struct {
	Url      string `koanf:"url"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
}

func InitConfig[T any](configFile iofs.FS) (T, error) {
	var config T
	k := koanf.New(".")
	configProvider := fs.Provider(configFile, "config.yaml")
	if err := k.Load(configProvider, yaml.Parser()); err != nil {
		return config, errors.New("cannot read config from file")
	}
	if err := k.Load(
		env.ProviderWithValue(
			EnvPrefix, ".", func(key string, value string) (string, any) {
				newKey := string_helper.SnakeToCamel(
					strings.Replace(
						strings.ToLower(
							strings.TrimPrefix(key, EnvPrefix),
						), "__", ".", -1,
					),
				)
				if strings.Contains(value, ",") {
					return newKey, strings.Split(value, ",")
				}
				return newKey, value
			},
		), nil,
	); err != nil {
		return config, err
	}

	if err := k.Unmarshal("", &config); err != nil {
		return config, err
	}
	return config, nil
}
