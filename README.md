# A Simple DNS Proxy/Forwarder in Go

This is something I just whipped and is a very basic example of how you can use
[`github.com/miekg/dns`](https://github.com/miekg/dns) to make a very simple DNS
forwarder.

Why? In my case, we are implementing the ability to use a specific resolver list
for DNS propagation checks in the [Terraform ACME
Provider](https://www.terraform.io/docs/providers/acme/index.html) and need a
way to test that our set resolvers actually took.
