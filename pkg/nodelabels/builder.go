/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nodelabels

import (
	"k8s.io/kops/pkg/apis/kops"
	"k8s.io/kops/util/pkg/reflectutils"
)

const (
	RoleLabelName15        = "kubernetes.io/role"
	RoleLabelName16        = "kubernetes.io/role"
	RoleMasterLabelValue15 = "master"
	RoleNodeLabelValue15   = "node"

	RoleLabelMaster16 = "node-role.kubernetes.io/master"
	RoleLabelNode16   = "node-role.kubernetes.io/node"

	RoleLabelControlPlane20 = "node-role.kubernetes.io/control-plane"
)

// BuildNodeLabels returns the node labels for the specified instance group
// This moved from the kubelet to a central controller in kubernetes 1.16
func BuildNodeLabels(cluster *kops.Cluster, instanceGroup *kops.InstanceGroup) map[string]string {
	isControlPlane := instanceGroup.Spec.Role == kops.InstanceGroupRoleMaster

	// Merge KubeletConfig for NodeLabels
	c := &kops.KubeletConfigSpec{}
	if isControlPlane {
		reflectutils.JSONMergeStruct(c, cluster.Spec.MasterKubelet)
	} else {
		reflectutils.JSONMergeStruct(c, cluster.Spec.Kubelet)
	}

	if instanceGroup.Spec.Kubelet != nil {
		reflectutils.JSONMergeStruct(c, instanceGroup.Spec.Kubelet)
	}

	nodeLabels := c.NodeLabels

	if isControlPlane {
		if nodeLabels == nil {
			nodeLabels = make(map[string]string)
		}
		for label, value := range BuildMandatoryControlPlaneLabels() {
			nodeLabels[label] = value
		}
	} else {
		if nodeLabels == nil {
			nodeLabels = make(map[string]string)
		}
		nodeLabels[RoleLabelNode16] = ""
		nodeLabels[RoleLabelName15] = RoleNodeLabelValue15
	}

	for k, v := range instanceGroup.Spec.NodeLabels {
		if nodeLabels == nil {
			nodeLabels = make(map[string]string)
		}
		nodeLabels[k] = v
	}

	return nodeLabels
}

// BuildMandatoryControlPlaneLabels returns the list of labels all CP nodes must have
func BuildMandatoryControlPlaneLabels() map[string]string {
	nodeLabels := make(map[string]string)
	nodeLabels[RoleLabelMaster16] = ""
	nodeLabels[RoleLabelControlPlane20] = ""
	nodeLabels[RoleLabelName15] = RoleMasterLabelValue15
	nodeLabels["kops.k8s.io/kops-controller-pki"] = ""
	return nodeLabels
}
