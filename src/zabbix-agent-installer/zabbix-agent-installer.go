package main

import (
	"fmt"
	"os"
	"os/exec"
	"log"
)

func usage() chan int {
	fmt.Printf("Usage: %s --os <ubuntu|centos>\n",os.Args[0])
	fmt.Println("    --os: Your OS type e.g. Ubuntu or CentOS")
	fmt.Println("    The OS version must be Ubuntu 14.04 or CentOS 6/7")

	ch <- 0
	return ch
}


func updateUbuntuSource() {
	aptUpdate, updateErr := exec.Command("sudo", "apt-get", "update").Output()
	if updateErr != nil {
		log.Fatal(updateErr)
	}
	fmt.Printf("Updating apt source...\n%s" ,aptUpdate)
}

func onUbuntu() chan int {
	updateUbuntuSource()

	installWget, installWgetErr := exec.Command("sudo","apt-get", "-y", "install", "wget").Output()
	if installWgetErr != nil {
		log.Fatal(installWgetErr)
	}
	fmt.Printf("Installing wget...\n%s", installWget)

	downZabbix, downZabbixErr := exec.Command("wget", "http://repo.zabbix.com/zabbix/2.2/ubuntu/pool/main/z/zabbix-release/zabbix-release_2.2-1+trusty_all.deb").Output()
	if downZabbixErr != nil {
		log.Fatal(downZabbixErr)
	}
	fmt.Printf("Downloing Zabbix repo...\n%s", downZabbix)

	installZabbixRepo, installZabbixRepoErr := exec.Command("sudo", "dpkg", "-i", "zabbix-release_2.2-1+trusty_all.deb").Output()
	if installZabbixRepoErr != nil {
		log.Fatal(installZabbixRepoErr)
	}
	fmt.Printf("Installing Zabbix repo...\n%s", installZabbixRepo)

	updateUbuntuSource()

	installZabbixAgent, installZabbixAgentErr := exec.Command("sudo", "apt-get", "-y", "install", "zabbix-agent").Output()
	if installZabbixAgentErr != nil {
		log.Fatal(installZabbixAgentErr)
	}
	fmt.Printf("Installing Zabbix Agent\n%s", installZabbixAgent)

	ch <- 0
	return ch
}

func onCentos() chan int {
	installZabbixRepo, installZabbixRepoErr := exec.Command("sudo", "rpm", "-ivh", "http://repo.zabbix.com/zabbix/2.2/rhel/6/x86_64/zabbix-release-2.2-1.el6.noarch.rpm").Output()
	if installZabbixRepoErr != nil {
		log.Fatal(installZabbixRepoErr)
	}
	fmt.Printf("Installing Zabbix repo...\n%s", installZabbixRepo)

	installZabbixAgent, installZabbixAgentErr := exec.Command("sudo", "yum", "-y", "install", "zabbix-agent").Output()
	if installZabbixAgentErr != nil {
		log.Fatal(installZabbixAgentErr)
	}
	fmt.Printf("Installing Zabbix Agent..\n%s", installZabbixAgent)

	ch <- 0
	return ch
}

var ch = make(chan int, 2)

func main() {
	if len(os.Args) != 3 {
		go usage()
		<- ch
	} else if os.Args[1] == "--os" && os.Args[2] == "ubuntu" {
		fmt.Println("Deploying Zabbix on Ubuntu...")
		go onUbuntu()
		<- ch
	} else if os.Args[1] == "--os" && os.Args[2] == "centos" {
		fmt.Println("Deploying Zabbix on CentOS")
		go onCentos()
		<- ch
	} else {
		go usage()
		<- ch
	}
}
