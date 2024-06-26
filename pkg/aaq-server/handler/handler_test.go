package handler

import (
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionv1 "k8s.io/api/admission/v1"
	authenticationv1 "k8s.io/api/authentication/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"kubevirt.io/application-aware-quota/pkg/client"
	"kubevirt.io/application-aware-quota/pkg/util"
	"kubevirt.io/application-aware-quota/pkg/util/patch"
	aaqv1 "kubevirt.io/application-aware-quota/staging/src/kubevirt.io/application-aware-quota-api/pkg/apis/core/v1alpha1"
	"kubevirt.io/application-aware-quota/tests/builders"
	"net/http"
)

var _ = Describe("Test handler of aaq server", func() {
	It("Pod should be gated", func() {
		pod := &v1.Pod{}
		podBytes, err := json.Marshal(pod)
		Expect(err).ToNot(HaveOccurred())

		v := Handler{
			request: &admissionv1.AdmissionRequest{
				Kind: metav1.GroupVersionKind{
					Kind: "Pod",
				},
				Object: runtime.RawExtension{
					Raw:    podBytes,
					Object: pod,
				},
				Operation: admissionv1.Create,
			},
		}
		admissionReview, err := v.Handle()
		Expect(err).ToNot(HaveOccurred())
		Expect(admissionReview.Response.Allowed).To(BeTrue())
		Expect(admissionReview.Response.Result.Code).To(Equal(int32(http.StatusAccepted)))
		Expect(admissionReview.Response.Result.Message).To(Equal(allowPodRequest))
	})

	DescribeTable("Pod Ungating should", func(ourController bool) {
		pod := &v1.Pod{
			Spec: v1.PodSpec{},
		}
		podBytes, err := json.Marshal(pod)
		Expect(err).ToNot(HaveOccurred())

		oldPod := &v1.Pod{
			Spec: v1.PodSpec{
				SchedulingGates: []v1.PodSchedulingGate{
					{
						util.AAQGate,
					},
				},
			},
		}
		oldPodBytes, err := json.Marshal(oldPod)
		Expect(err).ToNot(HaveOccurred())
		userName := util.ControllerServiceAccountName
		massage := aaqControllerPodUpdate
		allowed := true
		code := int32(http.StatusAccepted)
		if !ourController {
			userName = "some-controller"
			massage = invalidPodUpdate
			allowed = false
			code = int32(http.StatusForbidden)

		}
		v := Handler{
			request: &admissionv1.AdmissionRequest{
				Kind: metav1.GroupVersionKind{
					Kind: "Pod",
				},
				Object: runtime.RawExtension{
					Raw:    podBytes,
					Object: pod,
				},
				OldObject: runtime.RawExtension{
					Raw:    oldPodBytes,
					Object: oldPod,
				},
				Operation: admissionv1.Update,
				UserInfo: authenticationv1.UserInfo{
					Username: fmt.Sprintf("system:serviceaccount:%v:%v", util.DefaultAaqNs, userName),
				},
			},
			aaqNS: util.DefaultAaqNs,
		}
		admissionReview, err := v.Handle()
		Expect(err).ToNot(HaveOccurred())
		Expect(admissionReview.Response.Allowed).To(Equal(allowed))
		Expect(admissionReview.Response.Result.Code).To(Equal(code))
		Expect(admissionReview.Response.Result.Message).To(Equal(massage))
	},
		Entry(" be allowed for aaq-controller", true),
		Entry(" not be allowed for external users", false),
	)

	It("Pod Update without gate modifying our gate should be allowed", func() {
		pod := &v1.Pod{
			Spec: v1.PodSpec{
				SchedulingGates: []v1.PodSchedulingGate{
					{
						util.AAQGate,
					},
				},
			},
		}
		podBytes, err := json.Marshal(pod)
		Expect(err).ToNot(HaveOccurred())

		oldPod := &v1.Pod{
			Spec: v1.PodSpec{
				SchedulingGates: []v1.PodSchedulingGate{
					{
						util.AAQGate,
					},
					{
						"FakeGate",
					},
				},
			},
		}
		oldPodBytes, err := json.Marshal(oldPod)
		Expect(err).ToNot(HaveOccurred())

		v := Handler{
			request: &admissionv1.AdmissionRequest{
				Kind: metav1.GroupVersionKind{
					Kind: "Pod",
				},
				Object: runtime.RawExtension{
					Raw:    podBytes,
					Object: pod,
				},
				OldObject: runtime.RawExtension{
					Raw:    oldPodBytes,
					Object: oldPod,
				},
				Operation: admissionv1.Update,
				UserInfo: authenticationv1.UserInfo{
					Username: fmt.Sprintf("system:serviceaccount:%v:%v", "SomeNs", "someUser"),
				},
			},
			aaqNS: util.DefaultAaqNs,
		}
		admissionReview, err := v.Handle()
		Expect(err).ToNot(HaveOccurred())
		Expect(admissionReview.Response.Allowed).To(BeTrue())
		Expect(admissionReview.Response.Result.Code).To(Equal(int32(http.StatusAccepted)))
		Expect(admissionReview.Response.Result.Message).To(Equal(validPodUpdate))
	})

	DescribeTable("ARQ", func(operation admissionv1.Operation) {
		arq := builders.NewArqBuilder().WithNamespace("testNS").WithResource(v1.ResourceRequestsMemory, resource.MustParse("4Gi")).Build()
		arqBytes, err := json.Marshal(arq)
		Expect(err).ToNot(HaveOccurred())
		fakek8sCli := fake.NewSimpleClientset()
		ctrl := gomock.NewController(GinkgoT())
		cli := client.NewMockAAQClient(ctrl)
		cli.EXPECT().CoreV1().Times(1).Return(fakek8sCli.CoreV1())
		v := Handler{
			request: &admissionv1.AdmissionRequest{
				Kind: metav1.GroupVersionKind{
					Kind: "ApplicationAwareResourceQuota",
				},
				Object: runtime.RawExtension{
					Raw:    arqBytes,
					Object: arq,
				},
				Operation: operation,
				UserInfo: authenticationv1.UserInfo{
					Username: fmt.Sprintf("system:serviceaccount:%v:%v", "SomeNs", "someUser"),
				},
			},
			aaqNS:  util.DefaultAaqNs,
			aaqCli: cli,
		}
		admissionReview, err := v.Handle()
		Expect(err).ToNot(HaveOccurred())
		Expect(admissionReview.Response.Allowed).To(BeTrue())
		Expect(admissionReview.Response.Result.Code).To(Equal(int32(http.StatusAccepted)))
		Expect(admissionReview.Response.Result.Message).To(Equal(allowArqRequest))
	},
		Entry(" valid Update should be allowed", admissionv1.Update),
		Entry(" valid Creation should be allowed", admissionv1.Create),
	)

	It("pod with a non-empty spec.nodeName set", func() {
		const nodeName = "test-node"

		pod := &v1.Pod{
			Spec: v1.PodSpec{
				NodeName: nodeName,
			},
		}

		podBytes, err := json.Marshal(pod)
		Expect(err).ToNot(HaveOccurred())

		v := Handler{
			request: &admissionv1.AdmissionRequest{
				Kind: metav1.GroupVersionKind{
					Kind: "Pod",
				},
				Object: runtime.RawExtension{
					Raw:    podBytes,
					Object: pod,
				},
				Operation: admissionv1.Create,
			},
		}
		admissionReview, err := v.Handle()
		Expect(err).ToNot(HaveOccurred())
		Expect(admissionReview.Response.Allowed).To(BeTrue())

		var patches []patch.PatchOperation
		err = json.Unmarshal(admissionReview.Response.Patch, &patches)
		Expect(err).ToNot(HaveOccurred())

		removeNodeNamePatchFound := false
		nodeAffinityPatchFound := false
		wildcardTolerationPatchFound := false

		for _, p := range patches {
			switch p.Path {

			case "/spec/nodeName":
				removeNodeNamePatchFound = true
				Expect(p.Op).To(Equal(patch.PatchReplaceOp), "Patch is expected to replace nodeName with an empty string")
				Expect(p.Value).To(BeEmpty())

			case "/spec/affinity":
				nodeAffinityPatchFound = true

				var affinity *v1.Affinity
				valueBytes, err := json.Marshal(p.Value)
				Expect(err).ShouldNot(HaveOccurred())
				err = json.Unmarshal(valueBytes, &affinity)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms).To(ContainElement(
					v1.NodeSelectorTerm{
						MatchFields: []v1.NodeSelectorRequirement{
							{
								Key:      "metadata.name",
								Operator: v1.NodeSelectorOpIn,
								Values:   []string{pod.Spec.NodeName},
							},
						},
					},
				),
				)

			case "/spec/tolerations":
				wildcardTolerationPatchFound = true

				var tolerations []v1.Toleration
				valueBytes, err := json.Marshal(p.Value)
				Expect(err).ShouldNot(HaveOccurred())
				err = json.Unmarshal(valueBytes, &tolerations)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(tolerations).To(ContainElement(
					v1.Toleration{
						Key:      "", // i.e. any key
						Operator: v1.TolerationOpExists,
						Effect:   v1.TaintEffectNoSchedule,
					},
				))
			}
		}

		Expect(removeNodeNamePatchFound).To(BeTrue(), "Patch to remove nodeName is not found")
		Expect(nodeAffinityPatchFound).To(BeTrue(), "Patch to add node affinity is not found")
		Expect(wildcardTolerationPatchFound).To(BeTrue(), "Patch to a wildcard NoSchedule toleration is not found")
	})

	It("Operations on AAQ", func() {
		aaq := &aaqv1.AAQ{
			TypeMeta: metav1.TypeMeta{
				Kind:       "AAQ",
				APIVersion: "aaqs.aaq.kubevirt.io",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-aaq",
			},
		}
		aaqBytes, err := json.Marshal(aaq)
		Expect(err).ToNot(HaveOccurred())

		v := Handler{
			request: &admissionv1.AdmissionRequest{
				Kind: metav1.GroupVersionKind{
					Kind: aaq.Kind,
				},
				Object: runtime.RawExtension{
					Raw:    aaqBytes,
					Object: aaq,
				},
				Operation: admissionv1.Create,
			},
		}

		admissionReview, err := v.Handle()
		Expect(err).ToNot(HaveOccurred())
		Expect(admissionReview.Response.Allowed).To(BeFalse())
		Expect(admissionReview.Response.Result.Code).To(Equal(int32(http.StatusForbidden)))
		Expect(admissionReview.Response.Result.Message).To(Equal(onlySingleAAQInstaceIsAllowed))

	})
})
