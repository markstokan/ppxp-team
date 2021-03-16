package manifest_test

import (
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/planitest"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Logging", func() {

	Describe("timestamp format", func() {
		var manifest planitest.Manifest

		var jobs []string

		BeforeEach(func() {
			jobs = []string{
				"loggregator_agent_windows",
				"loggr-forwarder-agent-windows",
				"loggr-syslog-agent-windows",
				"prom_scraper_windows",
			}
		})

		When("logging_format_timestamp is set to rfc3339", func() {
			BeforeEach(func() {
				var err error
				// this test relies on the fixtures/tas_metadata.yml
				// that fixture sets "..cf.properties.logging_timestamp_format": "rfc3339"
				manifest, err = product.RenderManifest(map[string]interface{}{})
				Expect(err).NotTo(HaveOccurred())
			})

			It("sets format to rfc3339 on the logging jobs", func() {
				ig := "windows_diego_cell"

				for _, jobName := range jobs {
					job, err := manifest.FindInstanceGroupJob(ig, jobName)
					Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("%s job was not found on %s", jobName, ig))

					loggingFormatTimestamp, err := job.Property("logging/format/timestamp")
					Expect(err).NotTo(HaveOccurred())
					Expect(loggingFormatTimestamp).To(Equal("rfc3339"), fmt.Sprintf("%s failed", jobName))
				}
			})
		})
	})
	Describe("loggregator agent", func() {
		It("sets defaults on the loggregator agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "loggregator_agent_windows")
			Expect(err).NotTo(HaveOccurred())

			v2Api, err := agent.Property("loggregator/use_v2_api")
			Expect(err).ToNot(HaveOccurred())
			Expect(v2Api).To(BeTrue())

			tlsProps, err := agent.Property("loggregator/tls")
			Expect(err).ToNot(HaveOccurred())
			Expect(tlsProps).To(HaveKey("ca_cert"))

			expectSecureMetrics(agent)

			d, err := loadDomain("../../properties/logging.yml", "loggregator_agent_metrics_tls")
			Expect(err).ToNot(HaveOccurred())

			metricsProps, err := agent.Property("metrics")
			Expect(err).ToNot(HaveOccurred())
			Expect(metricsProps).To(HaveKeyWithValue("server_name", d))

			tlsAgentProps, err := agent.Property("loggregator/tls/agent")
			Expect(err).ToNot(HaveOccurred())
			Expect(tlsAgentProps).To(HaveKey("cert"))
			Expect(tlsAgentProps).To(HaveKey("key"))

			By("disabling udp")
			udpDisabled, err := agent.Property("disable_udp")
			Expect(err).NotTo(HaveOccurred())
			Expect(udpDisabled).To(BeTrue())

			By("getting the grpc port")
			port, err := agent.Property("grpc_port")
			Expect(err).NotTo(HaveOccurred())
			Expect(port).To(Equal(3459))

			By("setting tags on the emitted metrics")
			tags, err := agent.Property("tags")
			Expect(err).NotTo(HaveOccurred())
			Expect(tags).To(HaveKeyWithValue("product", "VMware Tanzu Application Service for Windows"))
			Expect(tags).NotTo(HaveKey("product_version"))
			Expect(tags).To(HaveKeyWithValue("system_domain", Not(BeEmpty())))
		})

		It("is enabled by default", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "loggregator_agent_windows")
			Expect(err).NotTo(HaveOccurred())

			_, err = agent.Property("loggregator_agent/enabled")
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when placement tags are configured by the user", func() {
			It("sets the placement tags on the emitted metrics", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".windows_diego_cell.placement_tags": "tag1,tag2",
				})
				Expect(err).NotTo(HaveOccurred())

				agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "loggregator_agent_windows")
				Expect(err).NotTo(HaveOccurred())

				tags, err := agent.Property("tags")
				Expect(err).NotTo(HaveOccurred())
				Expect(tags).To(HaveKeyWithValue("placement_tag", "tag1,tag2"))
			})
		})
	})

	Describe("forwarder agent", func() {
		It("sets defaults on the loggregator agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "loggr-forwarder-agent-windows")
			Expect(err).NotTo(HaveOccurred())

			expectSecureMetrics(agent)

			d, err := loadDomain("../../properties/logging.yml", "forwarder_agent_metrics_tls")
			Expect(err).ToNot(HaveOccurred())

			metricsProps, err := agent.Property("metrics")
			Expect(err).ToNot(HaveOccurred())
			Expect(metricsProps).To(HaveKeyWithValue("server_name", d))

			By("getting the grpc port")
			port, err := agent.Property("port")
			Expect(err).NotTo(HaveOccurred())
			Expect(port).To(Equal(3458))

			By("setting tags on the emitted metrics")
			tags, err := agent.Property("tags")
			Expect(err).NotTo(HaveOccurred())
			Expect(tags).To(HaveKeyWithValue("product", "VMware Tanzu Application Service for Windows"))
			Expect(tags).NotTo(HaveKey("product_version"))
			Expect(tags).To(HaveKeyWithValue("system_domain", Not(BeEmpty())))
		})

		Context("when placement tags are configured by the user", func() {
			It("sets the placement tags on the emitted metrics", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".windows_diego_cell.placement_tags": "tag1,tag2",
				})
				Expect(err).NotTo(HaveOccurred())

				agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "loggr-forwarder-agent-windows")
				Expect(err).NotTo(HaveOccurred())

				tags, err := agent.Property("tags")
				Expect(err).NotTo(HaveOccurred())
				Expect(tags).To(HaveKeyWithValue("placement_tag", "tag1,tag2"))
			})
		})
	})

	Describe("syslog agent", func() {
		It("sets defaults on the syslog agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "loggr-syslog-agent-windows")
			Expect(err).NotTo(HaveOccurred())

			expectSecureMetrics(agent)

			d, err := loadDomain("../../properties/logging.yml", "syslog_agent_metrics_tls")
			Expect(err).ToNot(HaveOccurred())

			metricsProps, err := agent.Property("metrics")
			Expect(err).ToNot(HaveOccurred())
			Expect(metricsProps).To(HaveKeyWithValue("server_name", d))

			port, err := agent.Property("port")
			Expect(err).NotTo(HaveOccurred())
			Expect(port).To(Equal(3460))

			tlsProps, err := agent.Property("tls")
			Expect(err).ToNot(HaveOccurred())
			Expect(tlsProps).To(HaveKey("ca_cert"))
			Expect(tlsProps).To(HaveKey("cert"))
			Expect(tlsProps).To(HaveKey("key"))

			cacheTlsProps, err := agent.Property("cache/tls")
			Expect(err).ToNot(HaveOccurred())
			Expect(cacheTlsProps).To(HaveKey("ca_cert"))
			Expect(cacheTlsProps).To(HaveKey("cert"))
			Expect(cacheTlsProps).To(HaveKey("key"))
			Expect(cacheTlsProps).To(HaveKeyWithValue("cn", "binding-cache"))
		})
	})
})

func expectSecureMetrics(job planitest.Manifest) {
	metricsProps, err := job.Property("metrics")
	Expect(err).ToNot(HaveOccurred())
	Expect(metricsProps).To(HaveKey("ca_cert"))
	Expect(metricsProps).To(HaveKey("cert"))
	Expect(metricsProps).To(HaveKey("key"))
	Expect(metricsProps).To(HaveKey("server_name"))
}

func loadDomain(file, property string) (string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	var certs []certEntry
	err = yaml.Unmarshal(b, &certs)
	if err != nil {
		return "", err
	}

	for _, c := range certs {
		if c.Name == property {
			if d, ok := c.Default.(map[interface{}]interface{}); ok {
				if doms, ok := d["domains"].([]interface{}); ok {
					return fmt.Sprintf("%v", doms[0]), nil
				}
			}

			return "", fmt.Errorf("property %s in %s incorrect", property, file)
		}
	}

	return "", fmt.Errorf("property %s not found in %s", property, file)
}

type certEntry struct {
	Name    string      `yaml:"name"`
	Default interface{} `yaml:"default"`
}