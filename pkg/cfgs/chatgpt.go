package cfgs

type ChatGpt struct {
	EndPoint    Expander `yaml:"EndPoint"`
	APIKey      Expander `yaml:"APIKey"`
	AccessToken Expander `yaml:"AccessToken"`
	Model       Expander `yaml:"Model"`
	Proxy       Expander `yaml:"Proxy"`
}
