package main

import "flag"

type Configuration struct {
	Quiet       *bool
	Region      *string
	ClientPort  *int
	UsePublicIP *bool
	AddMember   *bool
	Schema      *string
	Output      *string
}

func NewConfiguration() *Configuration {
	args := &Configuration{
		Quiet:       flag.Bool("quiet", false, "Disable log output"),
		Region:      flag.String("region", "eu-west-1", "Region to initialize the script."),
		ClientPort:  flag.Int("client-port", 2379, "ETCD Cient port"),
		UsePublicIP: flag.Bool("public", false, "Use EC2 Public IP for client URLs if available"),
		AddMember:   flag.Bool("add-member", true, "Add this etcd member explicitly to the cluster"),
		Schema:      flag.String("schema", "http", "Schema to communicate to the cluster, currently only 'http' works"),
		Output:      flag.String("output", "args", "Output format. Available options: args, env"),
	}

	flag.Parse()
	validaSchema(args.Schema)

	return args

}

func validateOutput(o *string) {
	if *o != "args" && *o != "env" {
		panic("Supported output types are 'args' and 'env'")
	}
}

func validaSchema(s *string) {
	if *s != "http" {
		panic("Only http scheam is supported!")
	}
}
