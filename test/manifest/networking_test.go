package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Networking", func() {
	Describe("Container networking", func() {
		var (
			inputProperties map[string]interface{}
			instanceGroup   string
		)

		BeforeEach(func() {
			if productName == "ert" {
				instanceGroup = "diego_cell"
			} else {
				instanceGroup = "compute"
			}
		})
		Context("when Silk is enabled", func() {
			BeforeEach(func() {
				inputProperties = map[string]interface{}{
					".properties.container_networking_interface_plugin": "silk",
				}
			})

			It("configures the cni_config_dir", func() {
				manifest, err := product.RenderService.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(instanceGroup, "garden-cni")
				Expect(err).NotTo(HaveOccurred())

				cniConfigDir, err := job.Property("cni_config_dir")
				Expect(err).NotTo(HaveOccurred())

				Expect(cniConfigDir).To(Equal("/var/vcap/jobs/silk-cni/config/cni"))
			})

			It("configures the cni_plugin_dir", func() {
				manifest, err := product.RenderService.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(instanceGroup, "garden-cni")
				Expect(err).NotTo(HaveOccurred())

				cniPluginDir, err := job.Property("cni_plugin_dir")
				Expect(err).NotTo(HaveOccurred())

				Expect(cniPluginDir).To(Equal("/var/vcap/packages/silk-cni/bin"))
			})
		})

		Context("when External is enabled", func() {
			BeforeEach(func() {
				inputProperties = map[string]interface{}{
					".properties.container_networking_interface_plugin": "external",
				}
			})

			It("configures the cni_config_dir", func() {
				manifest, err := product.RenderService.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(instanceGroup, "garden-cni")
				Expect(err).NotTo(HaveOccurred())

				cniConfigDir, err := job.Property("cni_config_dir")
				Expect(err).NotTo(HaveOccurred())

				Expect(cniConfigDir).To(Equal("/var/vcap/jobs/cni/config/cni"))
			})

			It("configures the cni_plugin_dir", func() {
				manifest, err := product.RenderService.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(instanceGroup, "garden-cni")
				Expect(err).NotTo(HaveOccurred())
				cniPluginDir, err := job.Property("cni_plugin_dir")
				Expect(err).NotTo(HaveOccurred())

				Expect(cniPluginDir).To(Equal("/var/vcap/packages/cni/bin"))
			})
		})
	})
})
