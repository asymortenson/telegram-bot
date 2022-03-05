package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken string
	ApiKey        string
	Fee float64
	Wallet string
	AdminChannel int64
	ExchangeChannel int64
	RequestLink string `mapstructure:"request_link"`
	Db struct {
		Dsn          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Messages Messages
}

type Messages struct {
	Responses
}



type Responses struct {
	PutPublicPage             string `mapstructure:"put_public_page"`
	BackBtn                   string `mapstructure:"back_button"`
	PaidBtn                   string `mapstructure:"paid_button"`
	DeclinePaidBtn            string `mapstructure:"decline_paid_button"`
	ChoosePublicPageBtn       string `mapstructure:"choose_public_page_button"`
	CreateRequestBtn          string `mapstructure:"create_request_button"`
	ApprovePublicPageBtn      string `mapstructure:"approve_public_page_button"`
	DeclinePublicPageBtn      string `mapstructure:"decline_public_page_button"`
	AfterSubmittingPublicPage string `mapstructure:"after_submitting_public_page"`
	RejectPublicPage          string `mapstructure:"reject_public_page"`
	AfterPaymentResponse      string `mapstructure:"after_payment_response"`
	FailedPaymentResponse     string `mapstructure:"failed_payment_response"`
	SuccessPaymentResponse    string `mapstructure:"success_payment_response"`
	PaymentMessage            string `mapstructure:"payment_message"`
	AlreadyPaid 			  string `mapstructure:"already_paid"`
	Signature 			  	  string `mapstructure:"signature"`
}

func Init() (*Config, error) {
	var cfg Config
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseEnv(cfg *Config) error {

	if err := viper.BindEnv("token"); err != nil {
		return err
	}

	if err := viper.BindEnv("api_key"); err != nil {
		return err
	}

	if err := viper.BindEnv("wallet"); err != nil {
		return err
	}

	if err := viper.BindEnv("fee"); err != nil {
		return err
	}
	if err := viper.BindEnv("admin_channel"); err != nil {
		return err
	}

	if err := viper.BindEnv("request_link"); err != nil {
		return err
	}

	if err := viper.BindEnv("exchange_channel"); err != nil {
		return err
	}

	if err := viper.BindEnv("db_dsn"); err != nil {
		return err
	}

	if err := viper.BindEnv("db_max_open_conns"); err != nil {
		return err
	}

	if err := viper.BindEnv("db_max_idle_conns"); err != nil {
		return err
	}

	if err := viper.BindEnv("db_max_idle_time"); err != nil {
		return err
	}

	
	cfg.TelegramToken = viper.GetString("token")
	cfg.ApiKey = viper.GetString("api_key")
	cfg.Wallet = viper.GetString("wallet")
	cfg.ExchangeChannel = viper.GetInt64("exchange_channel")
	cfg.AdminChannel = viper.GetInt64("admin_channel")
	cfg.Fee = viper.GetFloat64("fee")
	cfg.Db.Dsn = viper.GetString("db_dsn")
	cfg.RequestLink = viper.GetString("request_link")
	cfg.Db.MaxOpenConns = viper.GetInt("db_max_open_conns")
	cfg.Db.MaxIdleConns = viper.GetInt("db_max_idle_conns")
	cfg.Db.MaxIdleTime = viper.GetString("db_max_idle_time")

	return nil
}