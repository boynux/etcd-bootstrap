package main

import "flag"

type Configuration struct {
	Quiet      *bool
	Region     *string
	ClientPort *int
}

func NewConfiguration() *Configuration {
	args := &Configuration{
		Quiet:      flag.Bool("quiet", false, "Disable log output"),
		Region:     flag.String("region", "eu-west-1", "Region to initialize the script."),
		ClientPort: flag.Int("client-port", 2379, "ETCD Cient port"),
	}

	flag.Parse()

	return args

}
