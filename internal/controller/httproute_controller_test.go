package controller

import (
	"fmt"
	"strings"

	"github.com/cloudflare/cloudflare-go/v2"
	"github.com/cloudflare/cloudflare-go/v2/zero_trust"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTTPRoute Controller", func() {
	Context("When reconciling a resource", func() {
		It("should successfully reconcile the resource", func() {
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})

	Context("sortIngressByPathSpecificity", func() {
		It("should sort paths by length with longer paths first", func() {
			ingress := []zero_trust.TunnelConfigurationUpdateParamsConfigIngress{
				{Path: cloudflare.String("/"), Hostname: cloudflare.String("example.com"), Service: cloudflare.String("service1")},
				{Path: cloudflare.String("/api/v1"), Hostname: cloudflare.String("example.com"), Service: cloudflare.String("service2")},
				{Path: cloudflare.String("/api"), Hostname: cloudflare.String("example.com"), Service: cloudflare.String("service3")},
			}

			sortIngressByPathSpecificity(ingress)

			pathsAfterSort := extractPathsForTesting(ingress)
			Expect(pathsAfterSort[0]).To(Equal("/api/v1"))
			Expect(pathsAfterSort[1]).To(Equal("/api"))
			Expect(pathsAfterSort[2]).To(Equal("/"))
		})

		It("should prioritize non-wildcard paths over wildcard paths of same base length", func() {
			ingress := []zero_trust.TunnelConfigurationUpdateParamsConfigIngress{
				{Path: cloudflare.String("/token/*"), Hostname: cloudflare.String("example.com"), Service: cloudflare.String("service1")},
				{Path: cloudflare.String("/token"), Hostname: cloudflare.String("example.com"), Service: cloudflare.String("service2")},
			}

			sortIngressByPathSpecificity(ingress)

			pathsAfterSort := extractPathsForTesting(ingress)
			Expect(pathsAfterSort[0]).To(Equal("/token"))
			Expect(pathsAfterSort[1]).To(Equal("/token/*"))
		})

		It("should handle complex path combinations correctly", func() {
			ingress := []zero_trust.TunnelConfigurationUpdateParamsConfigIngress{
				{Path: cloudflare.String("/"), Hostname: cloudflare.String("api.example.com"), Service: cloudflare.String("service1")},
				{Path: cloudflare.String("/token/*"), Hostname: cloudflare.String("api.example.com"), Service: cloudflare.String("service2")},
				{Path: cloudflare.String("/asset-price"), Hostname: cloudflare.String("api.example.com"), Service: cloudflare.String("service3")},
				{Path: cloudflare.String("/api/v1/users"), Hostname: cloudflare.String("api.example.com"), Service: cloudflare.String("service4")},
			}

			sortIngressByPathSpecificity(ingress)

			pathsAfterSort := extractPathsForTesting(ingress)
			Expect(pathsAfterSort[0]).To(Equal("/api/v1/users"))
			Expect(pathsAfterSort[1]).To(Equal("/asset-price"))
			Expect(pathsAfterSort[2]).To(Equal("/token/*"))
			Expect(pathsAfterSort[3]).To(Equal("/"))
		})

		It("should handle empty and nil paths gracefully", func() {
			ingress := []zero_trust.TunnelConfigurationUpdateParamsConfigIngress{
				{Path: cloudflare.String(""), Hostname: cloudflare.String("example.com"), Service: cloudflare.String("service1")},
				{Path: cloudflare.String("/api"), Hostname: cloudflare.String("example.com"), Service: cloudflare.String("service2")},
				{Hostname: cloudflare.String("example.com"), Service: cloudflare.String("service3")}, // nil path
			}

			sortIngressByPathSpecificity(ingress)

			pathsAfterSort := extractPathsForTesting(ingress)
			Expect(pathsAfterSort[0]).To(Equal("/api"))
		})

		It("should sort alphabetically for paths of same length and wildcard status", func() {
			ingress := []zero_trust.TunnelConfigurationUpdateParamsConfigIngress{
				{Path: cloudflare.String("/zebra"), Hostname: cloudflare.String("example.com"), Service: cloudflare.String("service1")},
				{Path: cloudflare.String("/alpha"), Hostname: cloudflare.String("example.com"), Service: cloudflare.String("service2")},
			}

			sortIngressByPathSpecificity(ingress)

			pathsAfterSort := extractPathsForTesting(ingress)
			Expect(pathsAfterSort[0]).To(Equal("/alpha"))
			Expect(pathsAfterSort[1]).To(Equal("/zebra"))
		})
	})
})

func extractPathsForTesting(ingress []zero_trust.TunnelConfigurationUpdateParamsConfigIngress) []string {
	paths := make([]string, len(ingress))
	for i, rule := range ingress {
		path := fmt.Sprintf("%v", rule.Path)
		paths[i] = strings.Trim(path, "\"")
	}
	return paths
}
