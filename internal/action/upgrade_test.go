/*
Copyright 2022 The Flux authors

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

package action

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
	helmaction "helm.sh/helm/v3/pkg/action"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v2 "github.com/fluxcd/tanka-controller/api/v2"
)

func Test_newUpgrade(t *testing.T) {
	t.Run("new upgrade", func(t *testing.T) {
		g := NewWithT(t)

		obj := &v2.HelmRelease{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "upgrade",
				Namespace: "upgrade-ns",
			},
			Spec: v2.HelmReleaseSpec{
				Timeout: &metav1.Duration{Duration: time.Minute},
				Upgrade: &v2.Upgrade{
					Timeout: &metav1.Duration{Duration: 10 * time.Second},
					Force:   true,
				},
			},
		}

		got := newUpgrade(&helmaction.Configuration{}, obj, nil)
		g.Expect(got).ToNot(BeNil())
		g.Expect(got.Namespace).To(Equal(obj.Namespace))
		g.Expect(got.Timeout).To(Equal(obj.Spec.Upgrade.Timeout.Duration))
		g.Expect(got.Force).To(Equal(obj.Spec.Upgrade.Force))
	})

	t.Run("timeout fallback", func(t *testing.T) {
		g := NewWithT(t)

		obj := &v2.HelmRelease{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "upgrade",
				Namespace: "upgrade-ns",
			},
			Spec: v2.HelmReleaseSpec{
				Timeout: &metav1.Duration{Duration: time.Minute},
			},
		}

		got := newUpgrade(&helmaction.Configuration{}, obj, nil)
		g.Expect(got).ToNot(BeNil())
		g.Expect(got.Namespace).To(Equal(obj.Namespace))
		g.Expect(got.Timeout).To(Equal(obj.Spec.Timeout.Duration))
	})

	t.Run("applies options", func(t *testing.T) {
		g := NewWithT(t)

		obj := &v2.HelmRelease{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "upgrade",
				Namespace: "upgrade-ns",
			},
			Spec: v2.HelmReleaseSpec{},
		}

		got := newUpgrade(&helmaction.Configuration{}, obj, []UpgradeOption{
			func(upgrade *helmaction.Upgrade) {
				upgrade.Install = true
			},
			func(upgrade *helmaction.Upgrade) {
				upgrade.DryRun = true
			},
		})
		g.Expect(got).ToNot(BeNil())
		g.Expect(got.Install).To(BeTrue())
		g.Expect(got.DryRun).To(BeTrue())
	})

	t.Run("disable take ownership", func(t *testing.T) {
		g := NewWithT(t)

		obj := &v2.HelmRelease{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "upgrade",
				Namespace: "upgrade-ns",
			},
			Spec: v2.HelmReleaseSpec{
				Upgrade: &v2.Upgrade{
					DisableTakeOwnership: true,
				},
			},
		}

		got := newUpgrade(&helmaction.Configuration{}, obj, nil)
		g.Expect(got).ToNot(BeNil())
		g.Expect(got.TakeOwnership).To(BeFalse())
	})
}
