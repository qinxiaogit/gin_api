package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"time"
)

type ConsulService struct {
	IP   string
	Port int
	Tag  []string
	Name string
}

func RegisterService(ca string,cs *ConsulService){
	consulConfig := api.DefaultConfig()
	consulConfig.Address = ca
	client,err := api.NewClient(consulConfig)
	if err!=nil{
		fmt.Printf("NewClient error \n%v",err)
		return
	}
	agent := client.Agent()
	interval := time.Duration(10)*time.Second
	deregister := time.Duration(1)*time.Minute

	reg := &api.AgentServiceRegistration{
		ID: fmt.Sprintf("%v-%v-%v", cs.Name, cs.IP, cs.Port), // 服务节点的名称
		Name:cs.Name,
		Tags:cs.Tag,
		Port: cs.Port,
		Address: cs.IP,
		Check: &api.AgentServiceCheck{
			Interval: interval.String(),
			GRPC:                           fmt.Sprintf("%v:%v/%v", cs.IP, cs.Port, cs.Name), // grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
			DeregisterCriticalServiceAfter: deregister.String(),
		},
	}
	fmt.Printf("registing to %v\n",ca)
	if err := agent.ServiceRegister(reg);err!=nil{
		fmt.Printf("service Register error\n%v",err)
		return
	}
}