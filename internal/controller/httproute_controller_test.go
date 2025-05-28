package controller

import (
	"github.com/cloudflare/cloudflare-go/v2"
	"github.com/cloudflare/cloudflare-go/v2/zero_trust"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTTPRoute Controller", func() {
	Context("sortIngressByPathSpecificity", func() {
		It("should sort paths by specificity with longer paths first", func() {
			ingress := []zero_trust.TunnelConfigurationUpdateParamsConfigIngress{
				{Path: cloudflare.String("/"), Hostname: cloudflare.String("example.com")},
				{Path: cloudflare.String("/api/v1/users"), Hostname: cloudflare.String("example.com")},
				{Path: cloudflare.String("/api"), Hostname: cloudflare.String("example.com")},
			}

			sortIngressByPathSpecificity(ingress)

			paths := []string{
				extractPathSafely(ingress[0]),
				extractPathSafely(ingress[1]),
				extractPathSafely(ingress[2]),
			}

			Expect(paths[0]).To(Equal("/api/v1/users"), "First path should be the longest")
			Expect(paths[1]).To(Equal("/api"), "Second path should be medium length")
			Expect(paths[2]).To(Equal("/"), "Last path should be the shortest")
		})

		It("should handle wildcard paths correctly", func() {
			ingress := []zero_trust.TunnelConfigurationUpdateParamsConfigIngress{
				{Path: cloudflare.String("/"), Hostname: cloudflare.String("example.com")},
				{Path: cloudflare.String("/token/*"), Hostname: cloudflare.String("example.com")},
				{Path: cloudflare.String("/api/*"), Hostname: cloudflare.String("example.com")},
			}

			sortIngressByPathSpecificity(ingress)

			paths := []string{
				extractPathSafely(ingress[0]),
				extractPathSafely(ingress[1]),
				extractPathSafely(ingress[2]),
			}

			Expect(paths[0]).To(Equal("/token/*"), "First path should be the token wildcard")
			Expect(paths[1]).To(Equal("/api/*"), "Second path should be the api wildcard")
			Expect(paths[2]).To(Equal("/"), "Last path should be the root")
		})

		It("should handle entries without paths", func() {
			ingress := []zero_trust.TunnelConfigurationUpdateParamsConfigIngress{
				{Path: cloudflare.String("/api"), Hostname: cloudflare.String("example.com")},
				{Hostname: cloudflare.String("example.com")}, // no path
			}

			sortIngressByPathSpecificity(ingress)

			Expect(extractPathSafely(ingress[0])).To(Equal("/api"))
			Expect(extractPathSafely(ingress[1])).To(Equal(""))
		})
	})

	Context("When reconciling a resource", func() {
		It("should successfully reconcile the resource", func() {
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})
