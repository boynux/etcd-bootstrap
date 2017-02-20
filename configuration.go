package main

import "flag"

type Configuration struct {
	Quiet      *bool
	Region     *string
	ClientPort *int
	AddMember  *bool
	Schema     *string
}

func NewConfiguration() *Configuration {
	args := &Configuration{
		Quiet:      flag.Bool("quiet", false, "Disable log output"),
		Region:     flag.String("region", "eu-west-1", "Region to initialize the script."),
		ClientPort: flag.Int("client-port", 2379, "ETCD Cient port"),
		AddMember:  flag.Bool("add-member", true, "Add this etcd member implicitly to the cluster"),
		Schema:     flag.String("schema", "http", "Schema to communicate to the cluster, currently only 'http' works"),
	}

	flag.Parse()
	validaSchema(args.Schema)

	return args

}

func validaSchema(s *string) {
	if *s != "http" {
		panic("Only http scheam is supported!")
	}
}
