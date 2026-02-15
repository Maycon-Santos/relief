// Package proxy gerencia proxy reverso e configuração de rede.
package proxy

// TraefikRoute representa uma rota no Traefik
type TraefikRoute struct {
	Name    string
	Rule    string
	Service string
}

// TraefikService representa um serviço no Traefik
type TraefikService struct {
	Name string
	URL  string
}

// TraefikConfig representa a configuração completa do Traefik
type TraefikConfig struct {
	HTTP HTTPConfig `yaml:"http"`
}

// HTTPConfig contém configurações HTTP do Traefik
type HTTPConfig struct {
	Routers  map[string]Router  `yaml:"routers"`
	Services map[string]Service `yaml:"services"`
}

// Router representa um router do Traefik
type Router struct {
	Rule    string `yaml:"rule"`
	Service string `yaml:"service"`
}

// Service representa um serviço do Traefik
type Service struct {
	LoadBalancer LoadBalancer `yaml:"loadBalancer"`
}

// LoadBalancer representa a configuração de load balancer
type LoadBalancer struct {
	Servers []Server `yaml:"servers"`
}

// Server representa um servidor backend
type Server struct {
	URL string `yaml:"url"`
}
