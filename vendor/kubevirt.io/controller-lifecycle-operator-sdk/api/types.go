package api

import (
	conditions "github.com/openshift/custom-resource-status/conditions/v1"
	corev1 "k8s.io/api/core/v1"
)

// Phase is the current phase of the deployment
type Phase string

const (
	// PhaseDeploying signals that the resources are being deployed
	PhaseDeploying Phase = "Deploying"

	// PhaseDeployed signals that the resources are successfully deployed
	PhaseDeployed Phase = "Deployed"

	// PhaseDeleting signals that the resources are being removed
	PhaseDeleting Phase = "Deleting"

	// PhaseDeleted signals that the resources are deleted
	PhaseDeleted Phase = "Deleted"

	// PhaseError signals that the deployment is in an error state
	PhaseError Phase = "Error"

	// PhaseUpgrading signals that the resources are being deployed
	PhaseUpgrading Phase = "Upgrading"

	// PhaseEmpty is an uninitialized phase
	PhaseEmpty Phase = ""
)

// Status represents status of a operator configuration resource; must be inlined in the operator configuration resource status
type Status struct {
	Phase Phase `json:"phase,omitempty"`
	// A list of current conditions of the resource
	Conditions []conditions.Condition `json:"conditions,omitempty" optional:"true"`
	// The version of the resource as defined by the operator
	OperatorVersion string `json:"operatorVersion,omitempty" optional:"true"`
	// The desired version of the resource
	TargetVersion string `json:"targetVersion,omitempty" optional:"true"`
	// The observed version of the resource
	ObservedVersion string `json:"observedVersion,omitempty" optional:"true"`
}

// NodePlacement describes node scheduling configuration.
// +k8s:openapi-gen=true
type NodePlacement struct {
	// nodeSelector is the node selector applied to the relevant kind of pods
	// It specifies a map of key-value pairs: for the pod to be eligible to run on a node,
	// the node must have each of the indicated key-value pairs as labels
	// (it can have additional labels as well).
	// See https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
	// +kubebuilder:validation:Optional
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// affinity enables pod affinity/anti-affinity placement expanding the types of constraints
	// that can be expressed with nodeSelector.
	// affinity is going to be applied to the relevant kind of pods in parallel with nodeSelector
	// See https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity
	// +kubebuilder:validation:Optional
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// tolerations is a list of tolerations applied to the relevant kind of pods
	// See https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/ for more info.
	// These are additional tolerations other than default ones.
	// +kubebuilder:validation:Optional
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
}

// DeepCopyInto is copying the receiver, writing into out. in must be non-nil.
func (in *Status) DeepCopyInto(out *Status) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]conditions.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodePlacement) DeepCopyInto(out *NodePlacement) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(corev1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]corev1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodePlacement.
func (in *NodePlacement) DeepCopy() *NodePlacement {
	if in == nil {
		return nil
	}
	out := new(NodePlacement)
	in.DeepCopyInto(out)
	return out
}

// SwaggerDoc provides documentation for NodePlacement
func (NodePlacement) SwaggerDoc() map[string]string {
	return map[string]string{
		"":             "NodePlacement describes  node scheduling configuration.",
		"nodeSelector": "nodeSelector is the node selector applied to the relevant kind of pods\nIt specifies a map of key-value pairs: for the pod to be eligible to run on a node,\nthe node must have each of the indicated key-value pairs as labels\n(it can have additional labels as well).\nSee https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector\n+kubebuilder:validation:Optional\n+optional",
		"affinity":     "affinity enables pod affinity/anti-affinity placement expanding the types of constraints\nthat can be expressed with nodeSelector.\naffinity is going to be applied to the relevant kind of pods in parallel with nodeSelector\nSee https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity\n+kubebuilder:validation:Optional\n+optional",
		"tolerations":  "tolerations is a list of tolerations applied to the relevant kind of pods\nSee https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/ for more info.\nThese are additional tolerations other than default ones.\n+kubebuilder:validation:Optional\n+optional",
	}
}
