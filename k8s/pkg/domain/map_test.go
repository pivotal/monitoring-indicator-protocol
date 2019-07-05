package domain_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pivotal/monitoring-indicator-protocol/k8s/pkg/apis/indicatordocument/v1alpha1"
	"github.com/pivotal/monitoring-indicator-protocol/k8s/pkg/domain"
	"github.com/pivotal/monitoring-indicator-protocol/pkg/indicator"
)

func TestMap(t *testing.T) {
	t.Run("works with a complete document", func(t *testing.T) {
		g := NewGomegaWithT(t)

		k8sDoc := v1alpha1.IndicatorDocument{
			ObjectMeta: v1.ObjectMeta{
				Labels: map[string]string{"level": "high"},
			},
			Spec: v1alpha1.IndicatorDocumentSpec{
				Product: v1alpha1.Product{
					Name:    "my-product",
					Version: "my-version",
				},
				Indicators: []v1alpha1.IndicatorSpec{{
					Name:   "my-indicator",
					Promql: "my_promql",
					Alert: v1alpha1.Alert{
						For:  "10m",
						Step: "1m",
					},
					Thresholds: []v1alpha1.Threshold{{
						Level:    "critical",
						Operator: "eq",
						Value:    float64(100),
					}},
					Documentation: map[string]string{"docs": "explained"},
					Presentation: v1alpha1.Presentation{
						ChartType:    indicator.BarChart,
						CurrentValue: false,
						Frequency:    5,
						Labels:       []string{"pod", "app"},
					},
				}},
				Layout: v1alpha1.Layout{
					Owner:       "me",
					Title:       "my awesome indicators",
					Description: "enough said",
					Sections: []v1alpha1.Section{{
						Title:       "my section",
						Description: "the only section",
						Indicators:  []string{"my-indicator"},
					}},
				},
			},
		}

		i := indicator.Indicator{
			Name:   "my-indicator",
			PromQL: "my_promql",
			Thresholds: []indicator.Threshold{{
				Level:    "critical",
				Operator: indicator.EqualTo,
				Value:    100,
			}},
			Alert: indicator.Alert{
				For:  "10m",
				Step: "1m",
			},
			Documentation: map[string]string{"docs": "explained"},
			Presentation: indicator.Presentation{
				ChartType:    indicator.BarChart,
				CurrentValue: false,
				Frequency:    5,
				Labels:       []string{"pod", "app"},
			},
		}

		domainDoc := indicator.Document{
			Product: indicator.Product{
				Name:    "my-product",
				Version: "my-version",
			},
			Metadata:   map[string]string{"level": "high"},
			Indicators: []indicator.Indicator{i},
			Layout: indicator.Layout{
				Title:       "my awesome indicators",
				Description: "enough said",
				Sections: []indicator.Section{{
					Title:       "my section",
					Description: "the only section",
					Indicators:  []string{i.Name},
				}},
				Owner: "me",
			},
		}
		g.Expect(domain.Map(&k8sDoc)).To(BeEquivalentTo(domainDoc))
	})

	t.Run("works with duplicate indicators in the layout", func(t *testing.T) {
		g := NewGomegaWithT(t)

		k8sDoc := v1alpha1.IndicatorDocument{
			Spec: v1alpha1.IndicatorDocumentSpec{
				Product: v1alpha1.Product{
					Name:    "my-product",
					Version: "my-version",
				},
				Indicators: []v1alpha1.IndicatorSpec{{
					Name:   "my-indicator",
					Promql: "my_promql",
				}},
				Layout: v1alpha1.Layout{
					Owner:       "me",
					Title:       "my awesome indicators",
					Description: "enough said",
					Sections: []v1alpha1.Section{{
						Title:       "my section",
						Description: "the only section",
						Indicators:  []string{"my-indicator", "my-indicator"},
					}},
				},
			},
		}

		i := indicator.Indicator{
			Name:         "my-indicator",
			PromQL:       "my_promql",
			Thresholds:   []indicator.Threshold{},
			Presentation: indicator.Presentation{},
		}

		domainDoc := indicator.Document{
			Product: indicator.Product{
				Name:    "my-product",
				Version: "my-version",
			},
			Indicators: []indicator.Indicator{i},
			Layout: indicator.Layout{
				Title:       "my awesome indicators",
				Description: "enough said",
				Sections: []indicator.Section{{
					Title:       "my section",
					Description: "the only section",
					Indicators:  []string{i.Name, i.Name},
				}},
				Owner: "me",
			},
		}
		g.Expect(domain.Map(&k8sDoc)).To(BeEquivalentTo(domainDoc))
	})

	t.Run("uses the 'Undefined' operator when appropriate", func(t *testing.T) {
		g := NewGomegaWithT(t)

		k8sDoc := v1alpha1.IndicatorDocument{
			Spec: v1alpha1.IndicatorDocumentSpec{
				Product: v1alpha1.Product{
					Name:    "my-product",
					Version: "my-version",
				},
				Indicators: []v1alpha1.IndicatorSpec{{
					Name:   "my-indicator",
					Promql: "my_promql",
					Thresholds: []v1alpha1.Threshold{{
						Level: "WARNING",
					}},
				}},
			},
		}

		g.Expect(domain.Map(&k8sDoc).Indicators[0].Thresholds[0].Operator).
			To(BeEquivalentTo(indicator.Undefined))
	})
}
