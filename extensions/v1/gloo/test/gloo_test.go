package test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/go-utils/manifesttestutils"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/service-mesh-hub/pkg/render"
	"github.com/solo-io/service-mesh-hub/test"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("gloo extension test", func() {

	const (
		superglooNamesapce = "sm-marketplace"
		namespace          = "gloo-system"
		name               = "gloo"
	)

	var (
		spec       *v1.ApplicationSpec
		versionMap map[string]*v1.VersionedApplicationSpec
	)

	BeforeEach(func() {
		spec = test.LoadApplicationSpec("../spec.yaml")
		versionMap = make(map[string]*v1.VersionedApplicationSpec)
		for _, version := range spec.Versions {
			versionMap[version.Version] = version
		}
	})

	Context("istio", func() {
		const (
			meshName = "istio"
		)

		var (
			version      *v1.VersionedApplicationSpec
			inputs       render.ValuesInputs
			testManifest TestManifest
			testInput    = func(flavorName string, layers []render.LayerInput) render.ValuesInputs {
				return render.ValuesInputs{
					Name:             name,
					Flavor:           test.GetFlavor(flavorName, version),
					Layers:           layers,
					InstallNamespace: namespace,
					MeshRef: core.ResourceRef{
						Name:      meshName,
						Namespace: namespace,
					},
					SpecDefinedValues: version.ValuesYaml,
				}
			}
		)
		Context("0.18.35", func() {
			Context("with packaged flavor", func() {
				BeforeEach(func() {
					version = versionMap["0.18.35"]
					layers := []render.LayerInput{{
						LayerId:  "custom-resources",
						OptionId: "create",
					}}
					inputs = testInput("vanilla", layers)

				})

				It("has the correct number of resources with apiserver enabled", func() {
					inputs.Params = map[string]string{"apiServer.enable": "true"}
					rendered, err := render.ComputeResourcesForApplication(context.TODO(), inputs, version)
					Expect(err).NotTo(HaveOccurred())
					testManifest = NewTestManifestWithResources(rendered)
					Expect(testManifest.NumResources()).To(Equal(36))
				})

				It("has the correct number of resources with apiserver disabled", func() {
					inputs.Params = map[string]string{"apiServer.enable": "false"}
					rendered, err := render.ComputeResourcesForApplication(context.TODO(), inputs, version)
					Expect(err).NotTo(HaveOccurred())
					testManifest = NewTestManifestWithResources(rendered)
					Expect(testManifest.NumResources()).To(Equal(31))
				})
			})
		})
		Context("with custom flavor", func() {
			BeforeEach(func() {
				version = versionMap["0.18.35"]
				inputs = render.ValuesInputs{
					Name: name,
					Flavor: &v1.Flavor{
						Name: "custom-flavor",
					},
					InstallNamespace: namespace,
					MeshRef: core.ResourceRef{
						Name:      meshName,
						Namespace: namespace,
					},
					UserDefinedValues: "gateway:\n  upgrade: true",
					SpecDefinedValues: version.ValuesYaml,
				}
				rendered, err := render.ComputeResourcesForApplication(context.TODO(), inputs, version)
				Expect(err).NotTo(HaveOccurred())
				testManifest = NewTestManifestWithResources(rendered)
			})

			It("has the correct number of resources with gateway upgrade enabled", func() {
				Expect(testManifest.NumResources()).To(Equal(37))
			})

			It("has a job with gateway upgrade enabled", func() {
				testManifest.Expect("Job", "gloo-system", "gateway-conversion").NotTo(BeNil())
			})
		})
	})

})
