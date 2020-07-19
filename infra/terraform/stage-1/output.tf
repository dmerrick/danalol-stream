output default_vpc_cidr {
  description = "The CIDR block of the entire VPC"
  value       = module.default_vpc.default_vpc_cidr_block
}

output default_vpc_id {
  description = "The VPC ID of the default VPC"
  value       = module.default_vpc.default_vpc_id
}
