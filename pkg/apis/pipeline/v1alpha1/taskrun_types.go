/*
Copyright 2018 The Knative Authors.

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

package v1alpha1

import (
	"github.com/knative/pkg/apis"
	duckv1alpha1 "github.com/knative/pkg/apis/duck/v1alpha1"
	"github.com/knative/pkg/webhook"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Check that TaskRun may be validated and defaulted.
var _ apis.Validatable = (*TaskRun)(nil)
var _ apis.Defaultable = (*TaskRun)(nil)

// Assert that TaskRun implements the GenericCRD interface.
var _ webhook.GenericCRD = (*TaskRun)(nil)

// TaskRunSpec defines the desired state of TaskRun
type TaskRunSpec struct {
	TaskRef TaskRef     `json:"taskRef"`
	Trigger TaskTrigger `json:"trigger"`
	// +optional
	Inputs TaskRunInputs `json:"inputs,omitempty"`
	// +optional
	Outputs TaskRunOutputs `json:"outputs,omitempty"`
	Results Results        `json:"results"`
	// +optional
	Generation int64 `json:"generation,omitempty"`
	// +optional
	ServiceAccount string `json:"serviceAccount"`
}

// TaskRunInputs holds the input values that this task was invoked with.
type TaskRunInputs struct {
	// +optional
	Resources []TaskRunResourceVersion `json:"resourcesVersion,omitempty"`
	// +optional
	Params []Param `json:"params,omitempty"`
}

// TaskRunOutputs holds the output values that this task was invoked with.
type TaskRunOutputs struct {
	// +optional
	Resources []TaskRunResourceVersion `json:"resourcesVersion,omitempty"`
	// +optional
	Params []Param `json:"params,omitempty"`
}

// TaskTrigger defines a webhook style trigger to start a TaskRun
type TaskTrigger struct {
	TriggerRef TaskTriggerRef `json:"triggerRef"`
}

// TaskTriggerType indicates the mechanism by which this TaskRun was created.
type TaskTriggerType string

const (
	// TaskTriggerTypeManual indicates that this TaskRun was invoked manually by a user.
	TaskTriggerTypeManual TaskTriggerType = "manual"

	// TaskTriggerTypePipelineRun indicates that this TaskRun was created by a controller
	// attempting to realize a PipelineRun. In this case the `name` will refer to the name
	// of the PipelineRun.
	TaskTriggerTypePipelineRun TaskTriggerType = "pipelineRun"
)

// TaskTriggerRef describes what triggered this Task to run. It could be triggered manually,
// or it may have been part of a PipelineRun in which case this ref would refer
// to the corresponding PipelineRun.
type TaskTriggerRef struct {
	Type TaskTriggerType `json:"type"`
	// +optional
	Name string `json:"name,omitempty"`
}

var taskRunCondSet = duckv1alpha1.NewBatchConditionSet()

// TaskRunStatus defines the observed state of TaskRun
type TaskRunStatus struct {
	Results Results `json:"results"`
	Steps []StepRun `json:"steps"`
	// Conditions describes the set of conditions of this build.
	Conditions duckv1alpha1.Conditions `json:"conditions,omitempty"`
}

// GetCondition returns the Condition matching the given type.
func (tr *TaskRunStatus) GetCondition(t duckv1alpha1.ConditionType) *duckv1alpha1.Condition {
	return taskRunCondSet.Manage(tr).GetCondition(t)
}
func (tr *TaskRunStatus) InitializeConditions() {
	taskRunCondSet.Manage(tr).InitializeConditions()
}

// SetCondition sets the condition, unsetting previous conditions with the same
// type as necessary.
func (tr *TaskRunStatus) SetCondition(newCond *duckv1alpha1.Condition) {
	if newCond != nil {
		taskRunCondSet.Manage(tr).SetCondition(*newCond)
	}
}

// StepRun reports the results of running a step in the Task. Each
// task has the potential to succeed or fail (based on the exit code)
// and produces logs.
type StepRun struct {
	Name     string `json:"name"`
	LogsURL  string `json:"logsURL"`
	ExitCode int    `json:"exitCode"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TaskRun is the Schema for the taskruns API
// +k8s:openapi-gen=true
type TaskRun struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	Spec TaskRunSpec `json:"spec,omitempty"`
	// +optional
	Status TaskRunStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TaskRunList contains a list of TaskRun
type TaskRunList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TaskRun `json:"items"`
}

func (tr *TaskRun) SetDefaults() {}

// GetBuildRef for task
func (tr *TaskRun) GetBuildRef() corev1.ObjectReference {
	return corev1.ObjectReference{
		APIVersion: "build.knative.dev/v1alpha1",
		Kind:       "Build",
		Namespace:  tr.Namespace,
		Name:       tr.Name,
	}
}
