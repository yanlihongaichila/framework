package main

import (
	"fmt"
	"github.com/google/uuid"
	capi "github.com/hashicorp/consul/api"
)

func main() {
	client, err := capi.NewClient(capi.DefaultConfig())
	if err != nil {
		panic(err)
	}

	err = client.Agent().ServiceRegister(&capi.AgentServiceRegistration{
		ID:      uuid.NewString(),
		Name:    "test",
		Tags:    []string{"GRPC"},
		Port:    3306,
		Address: "127.0.0.1",
	})

	if err != nil {
		fmt.Println(err)
		return
	}
}
