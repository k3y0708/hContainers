package types

type Container struct {
	Name       string
	PortPrefix string
	Instance   string
	ID         string
	Image      string
	Version    string
	Runner     string
	Status     string
}

func (c Container) GetFullImageName() string {
	return c.Image + ":" + c.Version
}

func (c Container) GetPort() string {
	return c.PortPrefix + c.Instance
}

type ContainerInstances struct {
	Image   string
	Version string
	Count   int
	Healthy int
}
