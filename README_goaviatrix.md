# goaviatrix - An Aviatrix SDK for Go

The goaviatrix subfolder can be used as a golang SDK for the Aviatrix REST API. It's not feature complete, and currently is only known to be used for the aviatrix Terraform provider.


Full API docs are available at [apidoc](https://s3-us-west-2.amazonaws.com/avx-apidoc/index.htm)


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
    "github.com/AviatrixSystems/terraform-provider-aviatrix/goaviatrix"
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
