package metrics

import (
	"github.com/kamontat/fthelper/metric/v4/src/aggregators"
	"github.com/kamontat/fthelper/metric/v4/src/collectors"
	"github.com/kamontat/fthelper/metric/v4/src/connection"
	"github.com/kamontat/fthelper/metric/v4/src/freqtrade"
	"github.com/kamontat/fthelper/shared/caches"
	"github.com/kamontat/fthelper/shared/commandline/commands"
	"github.com/prometheus/client_golang/prometheus"
)

var FTPerformance = collectors.NewMetrics(
	collectors.NewMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName("freqtrade", "perf", "daily"),
			"Profit calculate by balance from yesterday and today (update once a day).",
			freqtrade.SummaryLabel(),
			nil,
		), func(desc *prometheus.Desc, conn connection.Http, param *commands.ExecutorParameter) []prometheus.Metric {
			var connection = freqtrade.ToConnection(conn)

			var data = caches.Global.Get(freqtrade.CACHE_DAILY_PERFORMANCE_BALANCE)

			var balance = freqtrade.NewBalance(connection)
			var previous = freqtrade.EmptyBalance()
			if data.IsExist() {
				previous = data.Data.(*freqtrade.Balance)
			}

			var labels = freqtrade.NewSummary(connection, param.Cache)
			var value, ok = aggregators.PercentChange(previous.CryptoValue, balance.CryptoValue)
			if !ok {
				param.Logger.Info("skip 'perf_daily' because previous is not exist")
				return emptyMetrics
			}

			return []prometheus.Metric{prometheus.MustNewConstMetric(
				desc,
				prometheus.GaugeValue,
				value,
				labels...,
			)}
		},
	),
)