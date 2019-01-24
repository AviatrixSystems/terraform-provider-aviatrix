# go-aviatrix - An Aviatrix SDK for Go

This is a golang SDK for the Aviatrix REST API. It's not feature complete, and currently is only known to be used for Terraform's `aviatrix` provider.

The code here is a mirror of [this](https://github.com/AviatrixSystems/terraform-provider-aviatrix/tree/master/vendor/github.com/AviatrixSystems/go-aviatrix/) folder in repo [terraform-provider-aviatrix](https://github.com/AviatrixSystems/terraform-provider-aviatrix).

Full API docs are available at [apidoc](https://s3-us-west-2.amazonaws.com/avx-apidoc/index.htm)

Any pull requests here will be ignored. Pull requests if any should be made at [terraform-provider-aviatrix](https://github.com/AviatrixSystems/terraform-provider-aviatrix). Commits into the vendor folder are automatically pushed into this repository


## Dependencies

* [ajg/form](https://github.com/ajg/form)
* [davecgh/go-spew](https://github.com/davecgh/go-spew.git)

## Example

```go
package main

import (
	"fmt"
    "log"
    "crypto/tls"
    "net/http"
    "github.com/go-aviatrix/goaviatrix"
)

func main() {

    tr := &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
	client, err := goaviatrix.NewClient("rakesh", "mypassword", "13.126.166.7", &http.Client{Transport: tr})
	if err != nil {
		log.Fatal(err)
	}
	err = client.CreateGateway(&goaviatrix.Gateway{
		Action: "connect_container",
		CloudType: 1,
		AccountName: "devops1",
		GwName: "avtxgw1",
		VpcID: "vpc-0d7b3664",
		VpcRegion: "ap-south-1",
		VpcSize: "t2.micro",
		VpcNet: "avtxgw3_sub1~~10.3.0.0/24~~ap-south-1a",
		})
	if err!=nil {
		log.Fatal(err)
	}
	err = client.DeleteGateway(&goaviatrix.Gateway{
		CloudType: 1,
		GwName: "avtxgw1",
		})
	if err!=nil {
		log.Fatal(err)
	}
	err = client.CreateTunnel(&goaviatrix.Tunnel{
		VpcName1: "avtxgw1",
		VpcName2: "avtxgw2",
	})
	if err!=nil {
		log.Fatal(err)
	}
	tun, err := client.GetTunnel(&goaviatrix.Tunnel{
		VpcName1: "avtxgw1",
		VpcName2: "avtxgw2",
	})
	if err!=nil {
		log.Fatal(err)
	}

	fmt.Println(tun.VpcName1, tun.VpcName2)

	err = client.DeleteTunnel(&goaviatrix.Tunnel{
		VpcName1: "avtxgw1",
		VpcName2: "avtxgw2",
	})

	if err!=nil {
		log.Fatal(err)
	}
}

```
