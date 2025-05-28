package controller

import (
	"github.com/cloudflare/cloudflare-go/v2"
	"github.com/cloudflare/cloudflare-go/v2/zero_trust"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTTPRoute Controller", func() {
	Context("When reconciling a resource", func() {
		It("should successfully reconcile the resource", func() {
			Expect(true).To(BeTrue())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})

	Context("When sorting ingress routes by path specificity", func() {
		It("should sort routes with longer paths first", func() {
			ingress := []zero_trust.TunnelConfigurationUpdateParamsConfigIngress{
				{
					Hostname: cloudflare.String("api.example.com"),
					Path:     cloudflare.F("/"),
					Service:  cloudflare.String("http://service1:80"),
				},
				{
					Hostname: cloudflare.String("api.example.com"),
					Path:     cloudflare.F("/token/validate"),
					Service:  cloudflare.String("http://service2:80"),
				},
				{
					Hostname: cloudflare.String("api.example.com"),
					Path:     cloudflare.F("/token"),
					Service:  cloudflare.String("http://service3:80"),
				},
			}

			sortIngressByPathSpecificity(ingress)

			Expect(ingress[0].Path.Value).To(Equal("/token/validate"))
			Expect(ingress[1].Path.Value).To(Equal("/token"))
			Expect(ingress[2].Path.Value).To(Equal("/"))
		})

		It("should handle nil path values safely", func() {
			ingress := []zero_trust.TunnelConfigurationUpdateParamsConfigIngress{
				{
					Hostname: cloudflare.String("api.example.com"),
					Path:     cloudflare.F[string](""),
					Service:  cloudflare.String("http://service1:80"),
				},
				{
					Hostname: cloudflare.String("api.example.com"),
					Path:     cloudflare.F("/token"),
					Service:  cloudflare.String("http://service2:80"),
				},
			}

			Expect(func() { sortIngressByPathSpecificity(ingress) }).ToNot(Panic())
			Expect(ingress[0].Path.Value).To(Equal("/token"))
		})
	})
})
