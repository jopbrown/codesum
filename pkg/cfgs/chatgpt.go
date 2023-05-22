package cfgs

type ChatGpt struct {
	EndPoint    Expander `yaml:"EndPoint" json:"end_point,omitempty"`
	APIKey      Expander `yaml:"APIKey" json:"api_key,omitempty"`
	AccessToken Expander `yaml:"AccessToken" json:"access_token,omitempty"`
	Model       Expander `yaml:"Model" json:"model,omitempty"`
	Proxy       Expander `yaml:"Proxy" json:"proxy,omitempty"`
}
