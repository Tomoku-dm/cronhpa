/*
Copyright 2022.

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

package controllers

import (
	"context"
	"reflect"
	"time"

	cronhpav1 "cronhpa/api/v1"

	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type CronHPAEvent = string
type CronHPA cronhpav1.CronHPA

const annotationNameSkip = "cronhpa.tomoku.github.com/skip"

const (
	CronHPAEventCreated     CronHPAEvent = "Created"
	CronHPAEventUpdated     CronHPAEvent = "Updated"
	CronHPAEventScheduled   CronHPAEvent = "Scheduled"
	CronHPAEventUnscheduled CronHPAEvent = "Unscheduled"
	CronHPAEventSkipped     CronHPAEvent = "Skipped"
	CronHPAEventNone        CronHPAEvent = ""
)

func (cronhpa *CronHPA) NewHPA() (*autoscalingv2beta2.HorizontalPodAutoscaler, error) {
	template := cronhpa.Spec.Template.DeepCopy()
	hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HorizontalPodAutoscaler",
			APIVersion: autoscalingv2beta2.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cronhpa.Name,
			Namespace: cronhpa.Namespace,
		},
		Spec: template.Spec,
	}
	if template.Metadata != nil {
		hpa.ObjectMeta.Labels = template.Metadata.Labels
		hpa.ObjectMeta.Annotations = template.Metadata.Annotations
	}
	return hpa, nil
}

func (cronhpa *CronHPA) CreateHPA(ctx context.Context, currentTime time.Time, reconciler *CronHPAReconciler) error {
	logger := log.FromContext(ctx)

	logger.Info("Create or update HPA")

	newhpa, err := cronhpa.NewHPA()
	if err != nil {
		return err
	}
	if err := controllerutil.SetControllerReference(cronhpa.ToCompatible(), newhpa, reconciler.Client.Scheme()); err != nil {
		return err
	}

	event := ""
	msg := ""
	hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{}
	if err := reconciler.Get(ctx, cronhpa.ToNamespacedName(), hpa); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		if err := reconciler.Create(ctx, newhpa); err != nil {
			return err
		}
		logger.Info("Created an HPA successfully")
		event = CronHPAEventCreated
		msg = "Created HPA"
	} else {
		if hpa.Annotations[annotationNameSkip] == "true" {
			logger.Info("Skip updating an HPA by an annotation")
			event = CronHPAEventSkipped
			msg = "Skipped updating HPA by an annotation"
		} else if reflect.DeepEqual(hpa.Spec, newhpa.Spec) {
			logger.Info("Skip updating an HPA with no changes")
			event = CronHPAEventSkipped
			msg = "Skipped updating HPA with no changes"
		} else {
			patch := client.MergeFrom(hpa)
			if err := reconciler.Patch(ctx, newhpa, patch); err != nil {
				return err
			}
			logger.Info("Updated an HPA successfully")
			event = CronHPAEventUpdated
			msg = "Updated HPA"
		}
	}

	if event != "" {
		cronhpa.Status.LastCronTimestamp = &metav1.Time{
			Time: currentTime,
		}
		reconciler.Recorder.Event(cronhpa.ToCompatible(), corev1.EventTypeNormal, event, msg)
	}

	return nil
}

func (cronhpa *CronHPA) ToCompatible() *cronhpav1.CronHPA {
	return (*cronhpav1.CronHPA)(cronhpa)
}

func (cronhpa *CronHPA) ToNamespacedName() types.NamespacedName {
	return types.NamespacedName{Namespace: cronhpa.Namespace, Name: cronhpa.Name}
}
