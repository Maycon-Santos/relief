package proxy

type TraefikRoute struct {
	Name    string
	Rule    string
	Service string
}

type TraefikService struct {
	Name string
	URL  string
}

type TraefikConfig struct {
	HTTP HTTPConfig `yaml:"http"`
}

type HTTPConfig struct {
	Routers  map[string]Router  `yaml:"routers"`
	Services map[string]Service `yaml:"services"`
}

type Router struct {
	Rule    string `yaml:"rule"`
	Service string `yaml:"service"`
}

type Service struct {
	LoadBalancer LoadBalancer `yaml:"loadBalancer"`
}

type LoadBalancer struct {
	Servers []Server `yaml:"servers"`
}

type Server struct {
	URL string `yaml:"url"`
}
