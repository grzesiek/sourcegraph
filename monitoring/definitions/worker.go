package definitions

import (
	"github.com/sourcegraph/sourcegraph/monitoring/definitions/shared"
	"github.com/sourcegraph/sourcegraph/monitoring/monitoring"
)

func Worker() *monitoring.Container {
	return &monitoring.Container{
		Name:        "worker",
		Title:       "Worker",
		Description: "Manages background processes.",
		Groups: []monitoring.Group{
			{
				Title:  "Precise code intelligence commit graph updater",
				Hidden: true,
				Rows: []monitoring.Row{
					{
						{
							Name:              "codeintel_commit_graph_queue_size",
							Description:       "commit graph queue size",
							Query:             `max(src_dirty_repositories_total)`,
							Warning:           monitoring.Alert().GreaterOrEqual(100, nil),
							Panel:             monitoring.Panel().LegendFormat("repositories with stale commit graphs"),
							Owner:             monitoring.ObservableOwnerCodeIntel,
							PossibleSolutions: "none",
						},
						{
							Name:              "codeintel_commit_graph_queue_growth_rate",
							Description:       "commit graph queue growth rate over 30m",
							Query:             `sum(increase(src_dirty_repositories_total[30m])) / sum(increase(src_codeintel_commit_graph_updater_total[30m]))`,
							Warning:           monitoring.Alert().GreaterOrEqual(5, nil),
							Panel:             monitoring.Panel().LegendFormat("rate of (enqueued / processed)"),
							Owner:             monitoring.ObservableOwnerCodeIntel,
							PossibleSolutions: "none",
						},
						{
							Name:              "codeintel_commit_graph_updater_99th_percentile_duration",
							Description:       "99th percentile successful commit graph updater operation duration over 5m",
							Query:             `histogram_quantile(0.99, sum by (le)(rate(src_codeintel_commit_graph_updater_duration_seconds_bucket{job=~"worker"}[5m])))`,
							Warning:           monitoring.Alert().GreaterOrEqual(20, nil),
							Panel:             monitoring.Panel().LegendFormat("update").Unit(monitoring.Seconds),
							Owner:             monitoring.ObservableOwnerCodeIntel,
							PossibleSolutions: "none",
						},
						{
							Name:              "codeintel_commit_graph_updater_errors",
							Description:       "commit graph updater errors every 5m",
							Query:             `sum(increase(src_codeintel_commit_graph_updater_errors_total{job=~"worker"}[5m]))`,
							Warning:           monitoring.Alert().GreaterOrEqual(20, nil),
							Panel:             monitoring.Panel().LegendFormat("errors"),
							Owner:             monitoring.ObservableOwnerCodeIntel,
							PossibleSolutions: "none",
						},
					},
				},
			},
			{
				Title:  "Precise code intelligence janitor",
				Hidden: true,
				Rows: []monitoring.Row{
					{
						{
							Name:              "codeintel_janitor_errors",
							Description:       "janitor errors every 5m",
							Query:             `sum(increase(src_codeintel_background_errors_total{job=~"worker"}[5m]))`,
							Warning:           monitoring.Alert().GreaterOrEqual(20, nil),
							Panel:             monitoring.Panel().LegendFormat("errors"),
							Owner:             monitoring.ObservableOwnerCodeIntel,
							PossibleSolutions: "none",
						},
						{
							Name:           "codeintel_upload_records_removed",
							Description:    "upload records expired or deleted every 5m",
							Query:          `sum(increase(src_codeintel_background_upload_records_removed_total{job=~"worker"}[5m]))`,
							NoAlert:        true,
							Panel:          monitoring.Panel().LegendFormat("uploads removed"),
							Owner:          monitoring.ObservableOwnerCodeIntel,
							Interpretation: "none",
						},
						{
							Name:           "codeintel_index_records_removed",
							Description:    "index records expired or deleted every 5m",
							Query:          `sum(increase(src_codeintel_background_index_records_removed_total{job=~"worker"}[5m]))`,
							NoAlert:        true,
							Panel:          monitoring.Panel().LegendFormat("indexes removed"),
							Owner:          monitoring.ObservableOwnerCodeIntel,
							Interpretation: "none",
						},
						{
							Name:           "codeintel_lsif_data_removed",
							Description:    "data for unreferenced upload records removed every 5m",
							Query:          `sum(increase(src_codeintel_background_uploads_purged_total{job=~"worker"}[5m]))`,
							NoAlert:        true,
							Panel:          monitoring.Panel().LegendFormat("uploads purged"),
							Owner:          monitoring.ObservableOwnerCodeIntel,
							Interpretation: "none",
						},
					},
					{
						{
							Name:              "codeintel_background_upload_resets",
							Description:       "upload records re-queued (due to unresponsive worker) every 5m",
							Query:             `sum(increase(src_codeintel_background_upload_resets_total{job=~"worker"}[5m]))`,
							Warning:           monitoring.Alert().GreaterOrEqual(20, nil),
							Panel:             monitoring.Panel().LegendFormat("uploads"),
							Owner:             monitoring.ObservableOwnerCodeIntel,
							PossibleSolutions: "none",
						},
						{
							Name:              "codeintel_background_upload_reset_failures",
							Description:       "upload records errored due to repeated reset every 5m",
							Query:             `sum(increase(src_codeintel_background_upload_reset_failures_total{job=~"worker"}[5m]))`,
							Warning:           monitoring.Alert().GreaterOrEqual(20, nil),
							Panel:             monitoring.Panel().LegendFormat("uploads"),
							Owner:             monitoring.ObservableOwnerCodeIntel,
							PossibleSolutions: "none",
						},
						{
							Name:              "codeintel_background_index_resets",
							Description:       "index records re-queued (due to unresponsive indexer) every 5m",
							Query:             `sum(increase(src_codeintel_background_index_resets_total{job=~"worker"}[5m]))`,
							Warning:           monitoring.Alert().GreaterOrEqual(20, nil),
							Panel:             monitoring.Panel().LegendFormat("indexes"),
							Owner:             monitoring.ObservableOwnerCodeIntel,
							PossibleSolutions: "none",
						},
						{
							Name:              "codeintel_background_index_reset_failures",
							Description:       "index records errored due to repeated reset every 5m",
							Query:             `sum(increase(src_codeintel_background_index_reset_failures_total{job=~"worker"}[5m]))`,
							Warning:           monitoring.Alert().GreaterOrEqual(20, nil),
							Panel:             monitoring.Panel().LegendFormat("indexes"),
							Owner:             monitoring.ObservableOwnerCodeIntel,
							PossibleSolutions: "none",
						},
					},
				},
			},
			{
				Title:  "Internal service requests",
				Hidden: true,
				Rows: []monitoring.Row{
					{
						shared.FrontendInternalAPIErrorResponses("worker", monitoring.ObservableOwnerCodeIntel).Observable(),
					},
				},
			},
			{
				Title:  shared.TitleContainerMonitoring,
				Hidden: true,
				Rows: []monitoring.Row{
					{
						shared.ContainerCPUUsage("worker", monitoring.ObservableOwnerCodeIntel).Observable(),
						shared.ContainerMemoryUsage("worker", monitoring.ObservableOwnerCodeIntel).Observable(),
					},
					{
						shared.ContainerMissing("worker", monitoring.ObservableOwnerCodeIntel).Observable(),
					},
				},
			},
			{
				Title:  shared.TitleProvisioningIndicators,
				Hidden: true,
				Rows: []monitoring.Row{
					{
						shared.ProvisioningCPUUsageLongTerm("worker", monitoring.ObservableOwnerCodeIntel).Observable(),
						shared.ProvisioningMemoryUsageLongTerm("worker", monitoring.ObservableOwnerCodeIntel).Observable(),
					},
					{
						shared.ProvisioningCPUUsageShortTerm("worker", monitoring.ObservableOwnerCodeIntel).Observable(),
						shared.ProvisioningMemoryUsageShortTerm("worker", monitoring.ObservableOwnerCodeIntel).Observable(),
					},
				},
			},
			{
				Title:  shared.TitleGolangMonitoring,
				Hidden: true,
				Rows: []monitoring.Row{
					{
						shared.GoGoroutines("worker", monitoring.ObservableOwnerCodeIntel).Observable(),
						shared.GoGcDuration("worker", monitoring.ObservableOwnerCodeIntel).Observable(),
					},
				},
			},
			{
				Title:  shared.TitleKubernetesMonitoring,
				Hidden: true,
				Rows: []monitoring.Row{
					{
						shared.KubernetesPodsAvailable("worker", monitoring.ObservableOwnerCodeIntel).Observable(),
					},
				},
			},
		},
	}
}