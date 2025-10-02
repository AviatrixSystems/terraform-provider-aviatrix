---
subcategory: "Firewall Network"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firewall_instance"
description: |-
  Creates and deletes Aviatrix Firewall Instances
---

# aviatrix_firewall_instance

The **aviatrix_firewall_instance** resource allows the creation and management of Aviatrix Firewall Instances.

This resource is used in [Aviatrix FireNet](https://docs.aviatrix.com/HowTos/firewall_network_faq.html) and [Aviatrix Transit FireNet](https://docs.aviatrix.com/HowTos/transit_firenet_faq.html) solutions, in conjunction with other resources that may include, and are not limited to: **firenet**, **firewall_instance_association**, **aws_tgw** and **transit_gateway** resources.

## Example Usage

```hcl
# Create an Aviatrix Firewall Instance
resource "aviatrix_firewall_instance" "test_firewall_instance" {
  vpc_id            = "vpc-032005cc371"
  firenet_gw_name   = "avx-firenet-gw"
  firewall_name     = "avx-firewall-instance"
  firewall_image    = "Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1"
  firewall_size     = "m5.xlarge"
  management_subnet = "10.4.0.16/28"
  egress_subnet     = "10.4.0.32/28"
}
```
```hcl
# Create an Aviatrix Firewall Instance with Native GWLB Enabled VPC
resource "aviatrix_firewall_instance" "test_firewall_instance" {
  vpc_id            = "vpc-032005cc371"
  firewall_name     = "avx-firewall-instance"
  firewall_image    = "Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1"
  firewall_size     = "m5.xlarge"
  management_subnet = "10.4.0.16/28"
  egress_subnet     = "10.4.0.32/28"
  zone              = "us-east-1a"
}
```
```hcl
# Create an Aviatrix Firewall Instance on GCP
resource "aviatrix_firewall_instance" "test_firewall_instance" {
  vpc_id                 = format("%s~-~%s", aviatrix_transit_gateway.test_transit_gateway.vpc_id, aviatrix_account.gcp.gcloud_project_id)
  firenet_gw_name        = aviatrix_transit_gateway.test_transit_gateway.gw_name
  firewall_name          = "gcp-firewall-instance"
  firewall_image         = "Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1"
  firewall_image_version = "9.0.9"
  firewall_size          = "n1-standard-4"
  management_vpc_id      = aviatrix_vpc.management_vpc.vpc_id
  management_subnet      = format("%s~~%s~~%s", aviatrix_vpc.management_vpc.subnets[0].cidr, aviatrix_vpc.management_vpc.subnets[0].region, aviatrix_vpc.management_vpc.subnets[0].name)
  egress_vpc_id          = aviatrix_vpc.egress_vpc.vpc_id
  egress_subnet          = format("%s~~%s~~%s", aviatrix_vpc.egress_vpc.subnets[0].cidr, aviatrix_vpc.egress_vpc.subnets[0].region, aviatrix_vpc.egress_vpc.subnets[0].name)
  zone                   = aviatrix_transit_gateway.test_transit_gateway.vpc_reg
}
```

## Argument Reference

The following arguments are supported:

### Required
~> As of Provider version R2.21.2+, the `vpc_id` of an OCI VCN has been changed from its name to its OCID.
* `vpc_id` - (Required) VPC ID of the Security VPC. For GCP, `vpc_id` must be in the form vpc_id~-~gcloud_project_id.
* `firenet_gw_name` - (Optional) Name of the primary FireNet gateway. **Required for all FireNet deployments that do not utilize the TGW-Integrated FireNet with AWS Native GWLB VPC.**
* `firewall_name` - (Required) Name of the firewall instance to be created.
* `firewall_image` - (Required) One of the AWS/Azure/GCP AMIs from various vendors such as Palo Alto Networks.
* `firewall_image_id` - (Optional) Firewall image ID. Only applicable to AWS. Please use AMI ID. Available as of provider version R2.19+.
* `firewall_size` - (Required) Instance size of the firewall. Example: "m5.xlarge".
* `management_subnet` - (Optional) Management Interface Subnet. Select the subnet whose name contains “gateway and firewall management”. For GCP, `management_subnet` must be in the form `cidr~~region~~name`. Required for Palo Alto Networks VM-Series and OCI Check Point firewalls. Otherwise, it must be empty.
* `egress_subnet` - (Required) Egress Interface Subnet. Select the subnet whose name contains “FW-ingress-egress”. For GCP, `egress_subnet` must be in the form `cidr~~region~~name`.
* `firewall_image_version` - (Optional) Version of firewall image. If not specified, Controller will automatically select the latest version available.
* `zone` - (Optional) Availability Zone. Required if creating a Firewall Instance with a Native AWS GWLB-enabled VPC. Applicable to AWS, Azure, and GCP only. Available as of provider version R2.17+.
* `management_vpc_id` - (Optional) Management VPC ID. Only used for GCP firewall. Required for Palo Alto Networks VM-Series, and required to be empty for Check Point or Fortinet series. Available as of provider version R2.18.1+.
* `egress_vpc_id` - (Optional) Egress VPC ID. Required for GCP. Available as of provider version R2.18.1+.
* `availability_domain` - (Optional) Availability domain. Required and valid only for OCI. Available as of provider version R2.19.3.
* `fault_domain` - (Optional) Fault domain. Required and valid only for OCI. Available as of provider version R2.19.3.

-> **NOTE:** Please use the data source `aviatrix_firewall_instance_images` to get the information for `firewall_image`, `firewall_image_version` and `firewall_size`.

Valid `firewall_image` values:

**AWS**
1. Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1
  - 10.2.14
  - 11.0.4-h6
  - 11.1.2-h3
  - 11.1.4-h7
  - 11.1.4-h13
  - 11.1.6-h7
  - 11.1.6-h10
  - 11.1.6-h14
  - 11.2.3-h3
  - 11.2.5
2. Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1 [VM-300]
  - 8.1.25-h1
3. Palo Alto Networks VM-Series Next-Generation Firewall Bundle 2
  - 10.2.10-h14
  - 10.2.14
  - 11.0.4-h6
  - 11.1.4-h7
  - 11.1.4-h13
  - 11.1.6-h7
  - 11.1.6-h10
  - 11.1.6-h14
  - 11.2.3-h3
  - 11.2.5
4. Palo Alto Networks VM-Series Next-Generation Firewall (BYOL)
  - 11.0.4-h6
  - 11.1.4-h7
  - 11.1.4-h13
  - 11.1.6-h7
  - 11.1.6-h10
  - 11.1.6-h14
  - 11.2.3-h3
  - 11.2.5
  - 11.2.8
  - 12.1.2
5. CloudGuard Network Security with Threat Prevention & SandBlast BYOL
  - R81-392.1729
  - R81.10-335.1727
  - R81.10-335.1788
  - R81.10-335.1878
  - R81.20-631.1753
  - R81.20-631.1856
  - R81.20-631.1857
  - R82-777.1836
  - R82-777.1869
6. CloudGuard Network Security Next-Gen Firewall with Threat Prevention
  - R81-392.1729
  - R81.10-335.1727
  - R81.10-335.1788
  - R81.10-335.1878
  - R81.20-631.1753
  - R81.20-631.1856
  - R81.20-631.1857
  - R82-777.1836
  - R82-777.1869
7. CloudGuard Network Security All-In-One
  - R81-392.1734
  - R81.10-335.1725
  - R81.10-335.1734
  - R81.10-335.1878
  - R81.20-634.1725
  - R81.20-634.1734
  - R81.20-634.1849
  - R82-777.1735
8. Fortinet FortiGate Next-Generation Firewall
  - (6.4.15)
  - (6.4.16)
  - (7.0.17)
  - (7.2.11)
  - (7.2.12)
  - (7.4.7)
  - (7.4.8)
  - (7.6.3)
  - (7.6.4)
9.  Fortinet FortiGate (BYOL) Next-Generation Firewall
  - (6.4.15)
  - (6.4.16)
  - (7.0.17)
  - (7.2.11)
  - (7.2.12)
  - (7.4.7)
  - (7.4.8)
  - (7.6.3)
  - (7.6.4)

**Azure**
1. Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1
  - 9.0.6
  - 9.0.9
  - 9.0.16
  - 9.1.0
  - 9.1.16
2. Palo Alto Networks VM-Series Next-Generation Firewall Bundle 2
  - 9.0.6
  - 9.0.9
  - 9.0.16
  - 9.1.0
  - 9.1.16
3. Palo Alto Networks VM-Series Next-Generation Firewall (BYOL)
  - 9.0.6
  - 9.0.9
  - 9.0.16
  - 9.1.0
  - 9.1.16
4. Palo Alto Networks VM-Series Flex Next-Generation Firewall Bundle 1
  - 11.1.604
  - 11.1.607
  - 11.2.5
  - 11.2.8
  - 11.2.303
5. Palo Alto Networks VM-Series Flex Next-Generation Firewall Bundle 2
  - 11.1.604
  - 11.1.607
  - 11.2.5
  - 11.2.8
  - 11.2.303
6. Palo Alto Networks VM-Series Flex Next-Generation Firewall (BYOL)
  - 11.1.607
  - 11.1.612
  - 11.2.5
  - 11.2.8
  - 11.2.303
7. Fortinet FortiGate (BYOL) Next-Generation Firewall
  - 7.4.5
  - 7.4.6
  - 7.4.7
  - 7.4.8
  - 7.6.0
8. Fortinet FortiGate (PAYG_2022) Next-Generation Firewall
  - 7.2.2
  - 7.2.3
  - 7.2.4
  - 7.2.5
  - 7.4.0
9. Fortinet FortiGate (PAYG_2023) Next-Generation Firewall
  - 7.4.5
  - 7.4.6
  - 7.4.7
  - 7.4.8
  - 7.6.0
10. Check Point CloudGuard IaaS Standalone (gateway + management)  R80.30 - Bring Your Own License
  - 8030.900200.0523
  - 8030.900200.0590
  - 8030.900200.0779
  - 8030.900200.0978
11. Check Point CloudGuard IaaS Single Gateway R80.30 - Pay As You Go (NGTP)
  - 8030.900273.0562
  - 8030.900273.0590
12. Check Point CloudGuard IaaS Single Gateway R80.30 - Pay As You Go (NGTX)
  - 8030.900273.0562
  - 8030.900273.0590
13. Check Point CloudGuard IaaS Standalone (gateway + management)  R80.40 - Bring Your Own License
  - 8040.900294.0586
  - 8040.900294.0752
  - 8040.900294.0978
  - 8040.900294.1108
  - 8040.900294.1255
  - 8040.900294.1266
  - 8040.900294.1300
  - 8040.900294.1335
  - 8040.900294.1370
  - 8040.900294.1383
  - 8040.900294.1475
  - 8040.900294.1526
  - 8040.900294.1560
14. Check Point CloudGuard IaaS Single Gateway R80.40 - Pay As You Go (NGTP)
  - 8040.900294.0593
  - 8040.900294.0640
  - 8040.900294.0682
  - 8040.900294.0728
  - 8040.900294.0801
  - 8040.900294.0907
  - 8040.900294.0978
  - 8040.900294.1108
  - 8040.900294.1234
  - 8040.900294.1255
  - 8040.900294.1266
  - 8040.900294.1300
  - 8040.900294.1331
  - 8040.900294.1370
  - 8040.900294.1383
  - 8040.900294.1475
  - 8040.900294.1498
  - 8040.900294.1526
  - 8040.900294.1564
15. Check Point CloudGuard IaaS Single Gateway R80.40 - Pay As You Go (NGTX)
  - 8040.900294.0593
  - 8040.900294.0640
  - 8040.900294.0682
  - 8040.900294.0728
  - 8040.900294.0801
  - 8040.900294.0907
  - 8040.900294.0978
  - 8040.900294.1108
  - 8040.900294.1234
  - 8040.900294.1255
  - 8040.900294.1266
  - 8040.900294.1300
  - 8040.900294.1331
  - 8040.900294.1370
  - 8040.900294.1383
  - 8040.900294.1475
  - 8040.900294.1498
  - 8040.900294.1526
  - 8040.900294.1564
16. Check Point CloudGuard IaaS Standalone (gateway + management)  R81 - Bring Your Own License
  - 8100.900392.0710
  - 8100.900392.0979
  - 8100.900392.1108
  - 8100.900392.1255
  - 8100.900392.1266
  - 8100.900392.1300
  - 8100.900392.1335
  - 8100.900392.1370
  - 8100.900392.1383
  - 8100.900392.1475
  - 8100.900392.1526
  - 8100.900392.1560
  - 8100.900392.1616
17. Check Point CloudGuard IaaS Single Gateway R81 - Pay As You Go (NGTP)
  - 8100.900392.0710
  - 8100.900392.0729
  - 8100.900392.0751
  - 8100.900392.0807
  - 8100.900392.0906
  - 8100.900392.0915
  - 8100.900392.0979
  - 8100.900392.1029
  - 8100.900392.1108
  - 8100.900392.1234
  - 8100.900392.1255
  - 8100.900392.1266
  - 8100.900392.1300
  - 8100.900392.1331
  - 8100.900392.1370
  - 8100.900392.1383
  - 8100.900392.1475
  - 8100.900392.1498
  - 8100.900392.1526
  - 8100.900392.1560
  - 8100.900392.1616
  - 8100.900392.1729
18. Check Point CloudGuard IaaS Single Gateway R81 - Pay As You Go (NGTX)
  - 8100.900392.0710
  - 8100.900392.0729
  - 8100.900392.0751
  - 8100.900392.0807
  - 8100.900392.0906
  - 8100.900392.0915
  - 8100.900392.0979
  - 8100.900392.1029
  - 8100.900392.1108
  - 8100.900392.1234
  - 8100.900392.1255
  - 8100.900392.1266
  - 8100.900392.1300
  - 8100.900392.1331
  - 8100.900392.1370
  - 8100.900392.1383
  - 8100.900392.1475
  - 8100.900392.1498
  - 8100.900392.1526
  - 8100.900392.1560
  - 8100.900392.1616
  - 8100.900392.1729
19. Check Point CloudGuard IaaS Standalone (gateway + management)  R81.10 - Bring Your Own License
  - 8110.900335.0883
  - 8110.900335.0979
  - 8110.900335.1108
  - 8110.900335.1226
  - 8110.900335.1255
  - 8110.900335.1266
  - 8110.900335.1300
  - 8110.900335.1335
  - 8110.900335.1370
  - 8110.900335.1383
  - 8110.900335.1475
  - 8110.900335.1526
  - 8110.900335.1560
  - 8110.900335.1612
  - 8110.900335.1648
  - 8110.900335.1681
  - 8110.900335.1734
  - 8110.900335.1878
20. Check Point CloudGuard IaaS Single Gateway R81.10 - Pay As You Go (NGTP)
  - 8110.900335.0883
  - 8110.900335.0979
  - 8110.900335.0985
  - 8110.900335.1108
  - 8110.900335.1234
  - 8110.900335.1255
  - 8110.900335.1266
  - 8110.900335.1300
  - 8110.900335.1331
  - 8110.900335.1366
  - 8110.900335.1370
  - 8110.900335.1383
  - 8110.900335.1475
  - 8110.900335.1498
  - 8110.900335.1526
  - 8110.900335.1546
  - 8110.900335.1560
  - 8110.900335.1612
  - 8110.900335.1648
  - 8110.900335.1681
  - 8110.900335.1709
  - 8110.900335.1715
  - 8110.900335.1727
  - 8110.900335.1736
  - 8110.900335.1788
  - 8110.900335.1878
  - 8110.900335.1895
  - 8110.900335.1903
21. Check Point CloudGuard IaaS Single Gateway R81.10 - Pay As You Go (NGTX)
  - 8110.900335.0883
  - 8110.900335.0979
  - 8110.900335.0985
  - 8110.900335.1108
  - 8110.900335.1234
  - 8110.900335.1255
  - 8110.900335.1266
  - 8110.900335.1300
  - 8110.900335.1331
  - 8110.900335.1366
  - 8110.900335.1370
  - 8110.900335.1383
  - 8110.900335.1475
  - 8110.900335.1498
  - 8110.900335.1526
  - 8110.900335.1546
  - 8110.900335.1560
  - 8110.900335.1612
  - 8110.900335.1648
  - 8110.900335.1681
  - 8110.900335.1709
  - 8110.900335.1715
  - 8110.900335.1727
  - 8110.900335.1736
  - 8110.900335.1788
  - 8110.900335.1878
  - 8110.900335.1895
  - 8110.900335.1903
22. Check Point CloudGuard IaaS Standalone (gateway + management)  R81.20 - Bring Your Own License
  - 8120.900631.01243
  - 8120.900631.01266
  - 8120.900631.01306
  - 8120.900631.01335
  - 8120.900631.01370
  - 8120.900631.01383
  - 8120.900631.01475
  - 8120.900631.01526
  - 8120.900631.01560
  - 8120.900634.01599
  - 8120.900634.01641
  - 8120.900634.01723
  - 8120.900634.01734
  - 8120.900634.01849
23. Check Point CloudGuard IaaS Single Gateway R81.20 - Pay As You Go (NGTP)
  - 8120.900627.01193
  - 8120.900631.01245
  - 8120.900631.01266
  - 8120.900631.01306
  - 8120.900631.01331
  - 8120.900631.01366
  - 8120.900631.01370
  - 8120.900631.01383
  - 8120.900631.01475
  - 8120.900631.01526
  - 8120.900631.01544
  - 8120.900631.01560
  - 8120.900634.01611
  - 8120.900634.01641
  - 8120.920631.01669
  - 8120.920631.01709
  - 8120.920631.01716
  - 8120.920631.01727
  - 8120.920631.01731
  - 8120.920631.01736
  - 8120.920631.01753
  - 8120.920631.01781
  - 8120.920631.01849
  - 8120.920631.01896
  - 8120.920631.01903
24. Check Point CloudGuard IaaS Single Gateway R81.20 - Pay As You Go (NGTX)
  - 8120.900627.01193
  - 8120.900631.01245
  - 8120.900631.01266
  - 8120.900631.01306
  - 8120.900631.01331
  - 8120.900631.01366
  - 8120.900631.01370
  - 8120.900631.01383
  - 8120.900631.01475
  - 8120.900631.01526
  - 8120.900631.01544
  - 8120.900631.01560
  - 8120.900634.01611
  - 8120.900634.01641
  - 8120.920631.01669
  - 8120.920631.01709
  - 8120.920631.01716
  - 8120.920631.01727
  - 8120.920631.01731
  - 8120.920631.01736
  - 8120.920631.01753
  - 8120.920631.01781
  - 8120.920631.01849
  - 8120.920631.01896
  - 8120.920631.01903
25. Check Point CloudGuard IaaS Standalone (gateway + management)  R82 - Bring Your Own License
  - 8200.900777.1695
  - 8200.900777.1734
26. Check Point CloudGuard IaaS Single Gateway R82 - Pay As You Go (NGTP)
  - 8200.777.1695
  - 8200.900777.1715
  - 8200.900777.1736
  - 8200.900777.1832
  - 8200.900777.1836
  - 8200.900777.1869
  - 8200.900777.1897
  - 8200.900777.1903
27. Check Point CloudGuard IaaS Single Gateway R82 - Pay As You Go (NGTX) 
  - 8200.900777.1695
  - 8200.900777.1715
  - 8200.900777.1736
  - 8200.900777.1832
  - 8200.900777.1836
  - 8200.900777.1869
  - 8200.900777.1897
  - 8200.900777.1903

**GCP**
1. Palo Alto Networks VM-Series Next-Generation Firewall BUNDLE1
  - 8.1.25-h1
  - 9.0.16-h5
2. Palo Alto Networks VM-Series Next-Generation Firewall BUNDLE2
  - 8.1.25-h1
  - 9.0.16-h5
3. Palo Alto Networks VM-Series Next-Generation Firewall BYOL
  - 9.0.16-h5
  - 8.1.25-h1
4. Palo Alto Networks VM-Series Flex Next-Generation Firewall BUNDLE1
  - 10.1.14-h6
  - 10.1.14-h8
  - 10.1.14-h9
  - 10.2.10-h9
  - 10.2.10-h12
  - 10.2.10-h14
  - 11.1.4-h13
  - 11.1.6-h7
  - 11.2.5
  - 11.2.6
5. Palo Alto Networks VM-Series Flex Next-Generation Firewall BUNDLE2
  - 10.1.14-h6
  - 10.1.14-h8
  - 10.1.14-h9
  - 10.2.10-h9
  - 10.2.10-h12
  - 10.2.10-h14
  - 11.1.4-h13
  - 11.1.6-h7
  - 11.2.5
  - 11.2.6
6. Palo Alto Networks VM-Series Flex Next-Generation Firewall BUNDLE3
  - 10.2.10-h9
  - 10.2.10-h12
  - 10.2.10-h14
  - 11.0.4-h6
  - 11.1.4-h7
  - 11.1.4-h13
  - 11.1.6-h7
  - 11.2.3-h3
  - 11.2.5
  - 11.2.6
7. Palo Alto Networks VM-Series Flex Next-Generation Firewall BYOL
  - 10.1.14-h8
  - 10.1.14-h9
  - 10.2.10-h12
  - 10.2.10-h14
  - 11.1.4-h13
  - 11.1.6-h7
  - 11.2.5
  - 11.2.6
  - 11.2.8
  - 12.1.2
8. Fortinet FortiGate Next-Generation Firewall
  - 6.4.1.6
  - 7.0.1.7
  - 7.2.1.1
  - 7.2.1.2
  - 7.4.7
  - 7.4.8
  - 7.6.2
  - 7.6.3
  - 7.6.4
9.  Fortinet FortiGate Next-Generation Firewall (BYOL)
  - 6.4.1.6
  - 7.0.1.7
  - 7.2.1.1
  - 7.2.1.2
  - 7.4.7
  - 7.4.8
  - 7.6.2
  - 7.6.3
  - 7.6.4
10. Check Point CloudGuard IaaS Firewall & Threat Prevention (Gateway only)
  - R80.40-294.991001526
  - R80.40-294.991001555
  - R80.40-294.991001560
  - R80.40-294.991001564
  - R81-392.991001498
  - R81-392.991001526
  - R81-392.991001555
  - R81-392.991001560
  - R81-392.991001616
  - R81-392.991001729
11. Check Point CloudGuard IaaS Firewall & Threat Prevention (Gateway only) (BYOL)
  - R80.40-294.991001526
  - R80.40-294.991001555
  - R80.40-294.991001560
  - R80.40-294.991001564
  - R81-392.991001498
  - R81-392.991001526
  - R81-392.991001555
  - R81-392.991001560
  - R81-392.991001616
  - R81-392.991001729
12. Check Point CloudGuard IaaS Firewall & Threat Prevention (Standalone)
  - R80.40-294.991001475
  - R80.40-294.991001526
  - R80.40-294.991001555
  - R80.40-294.991001560
  - R81-392.991001383
  - R81-392.991001475
  - R81-392.991001526
  - R81-392.991001555
  - R81-392.991001560
  - R81-392.991001616
13. Check Point CloudGuard IaaS Firewall & Threat Prevention (Standalone) (BYOL)
  - R80.40-294.991001526
  - R80.40-294.991001555
  - R80.40-294.991001560
  - R81-392.991001526
  - R81-392.991001555
  - R81-392.991001560
  - R81-392.991001616
  - R81.10-335.991001526
  - R81.10-335.991001555
  - R81.10-335.991001560

**OCI**
1. Palo Alto Networks VM-Series Bundle1 - 4 OCPUs
  - 10.1.14-h6
  - 10.1.14-h8
  - 10.1.14-h9
  - 10.2.10-h12
  - 10.2.10-h14
  - 11.1.4-h13
  - 11.1.6-h7
  - 11.2.5
  - 11.2.8
  - 12.1.2
2. Palo Alto Networks VM-Series Bundle1 - 8 OCPUs
  - 10.1.14-h6
  - 10.1.14-h8
  - 10.1.14-h9
  - 10.2.10-h12
  - 10.2.10-h14
  - 11.1.4-h13
  - 11.1.6-h7
  - 11.2.5
  - 11.2.8
  - 12.1.2
3. Palo Alto Networks VM-Series Bundle2 - 4 OCPUs
  - 10.1.14-h6
  - 10.1.14-h8
  - 10.1.14-h9
  - 10.2.10-h12
  - 10.2.10-h14
  - 11.1.4-h13
  - 11.1.6-h7
  - 11.2.5
  - 11.2.8
  - 12.1.2
4. Palo Alto Networks VM-Series Bundle2 - 8 OCPUs
  - 10.1.14-h6
  - 10.1.14-h8
  - 10.1.14-h9
  - 10.2.10-h12
  - 10.2.10-h14
  - 11.1.4-h13
  - 11.1.6-h7
  - 11.2.5
  - 11.2.8
  - 12.1.2
5. Palo Alto Networks VM-Series Bundle3 - 4 OCPUs
  - 10.2.10-h9
  - 10.2.10-h12
  - 10.2.10-h14
  - 11.0.4-h6
  - 11.1.3
  - 11.1.4-h7
  - 11.1.4-h13
  - 11.1.6-h7
  - 11.2.3-h3
  - 11.2.5
6. Palo Alto Networks VM-Series Bundle3 - 8 OCPUs
  - 10.2.10-h9
  - 10.2.10-h12
  - 10.2.10-h14
  - 11.0.4-h6
  - 11.1.3
  - 11.1.4-h7
  - 11.1.4-h13
  - 11.1.6-h7
  - 11.2.3-h3
  - 11.2.5
7. Palo Alto Networks VM-Series Next Generation Firewall
  - 10.1.14-h6
  - 10.1.14-h8
  - 10.1.14-h9
  - 10.2.10-h12
  - 10.2.10-h14
  - 11.1.4-h13
  - 11.1.6-h7
  - 11.2.5
  - 11.2.8
  - 12.1.2
8. CloudGuard Next-Gen Firewall w/ Threat Prevention - 4 OCPUs
  - R81.10_with_JHF_45
  - R81.10_with_JHF_150
  - R81.10_with_JHF_177
  - R81.20_with_JHF_99
  - R81.20_with_JHF_105_v2
  - R81.20_JHF53
  - R81.20_rev1.0
  - R81_with_JHF_65
  - R82_with_JHF_19
  - R82_with_JHF_34
9.  CloudGuard Next-Gen Firewall w/ Threat Prevention and SandBlast - 4 OCPUs
  - R81.10_with_JHF_45
  - R81.10_with_JHF_150
  - R81.10_with_JHF_177
  - R81.20_with_JHF_99
  - R81.20_with_JHF_105_v2
  - R81.20_JHF53
  - R81.20_rev1.0
  - R81_with_JHF_65
  - R82_with_JHF_19
  - R82_with_JHF_34
10. CloudGuard Next-Gen Firewall with Threat Prevention and SandBlast - BYOL
  - R81.10_with_JHF_150
  - R81.10_with_JHF_177
  - R81.10_with_JHF_45
  - R81.20_JHF53
  - R81.20_rev1.0
  - R81.20_with_JHF_105_v2
  - R81.20_with_JHF_99
  - R81_with_JHF_65
  - R82_with_JHF_19
  - R82_with_JHF_34

### Authentication method
* `key_name`- (Optional) Applicable to AWS deployment only. AWS Key Pair name. If not provided a Key Pair will be generated.
* `username`- (Optional) Applicable to Azure or AzureGov deployment only. "admin" as a username is not accepted.
* `password`- (Optional) Applicable to Azure or AzureGov deployment only.
* `ssh_public_key` - (Optional) Applicable to Azure or AzureGov deployment only.

### Advanced Options
* `iam_role` - (Optional) Only available for AWS. In advanced mode, create an IAM Role on the AWS account that launched the FireNet gateway. Create a policy to attach to the role. The policy is to allow access to "Bootstrap Bucket".
* `bootstrap_bucket_name`- (Optional) Only available for AWS and GCP. For GCP, only Palo Alto Networks VM-Series deployment can use this attribute. In advanced mode, specify a bootstrap bucket name where the initial configuration and policy file is stored.
* `bootstrap_storage_name` - (Optional) Advanced option. Bootstrap storage name. Applicable to Azure or AzureGov and Palo Alto Networks VM-Series/Fortinet Series deployment only. Available as of provider version R2.17.1+.
* `storage_access_key` - (Optional) Advanced option. Storage access key. Applicable to Azure or AzureGov and Palo Alto Networks VM-Series deployment only. Available as of provider version R2.17.1+.
* `file_share_folder` - (Optional) Advanced option. File share folder. Applicable to Azure or AzureGov and Palo Alto Networks VM-Series deployment only. Available as of provider version R2.17.1+.
* `share_directory` - (Optional) Advanced option. Share directory. Applicable to Azure or AzureGov and Palo Alto Networks VM-Series deployment only. Available as of provider version R2.17.1+.
* `sic_key` - (Optional) Advanced option. Sic key. Applicable to Check Point Series deployment only.
* `container_folder` - (Optional) Advanced option. Container folder. Applicable to Azure or AzureGov and Fortinet Series deployment only.
* `sas_url_config` - (Optional) Advanced option. SAS URL Config. Applicable to Azure or AzureGov and Fortinet Series deployment only.
* `sas_url_license` - (Optional) Advanced option. SAS URL License. Applicable to Azure or AzureGov and Fortinet Series deployment only.
* `user_data` - (Optional) Advanced option. User Data. Type: String.

### Misc.
* `tags` - (Optional) Mapping of key value pairs of tags for a firewall instance. Only available for AWS, AWSGov, GCP and Azure firewall instances. For AWS, AWSGov and Azure allowed characters are: letters, spaces, and numbers plus the following special characters: + - = . _ : @. For GCP allowed characters are: lowercase letters, numbers, "-" and "_". Example: {"key1" = "value1", "key2" = "value2"}.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `instance_id`- ID of the firewall instance created.
* `lan_interface`- ID of Lan Interface created.
* `management_interface`- ID of Management Interface created.
* `egress_interface`- ID of Egress Interface created.
* `public_ip`- Management Public IP.
* `cloud_type` - Cloud Type.
* `gcp_vpc_id` - GCP Only. The current VPC ID.

## Import

**firewall_instance** can be imported using the `instance_id`. For Azure or AzureGov FireNet instances, the value will be the `firewall_name` concatenated with a ":" and the Resource Group of the `vpc_id` set for that instance. e.g.

```
$ terraform import aviatrix_firewall_instance.test instance_id
```