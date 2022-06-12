package nats

type NatsAnalysisConfig struct {
	consumerCalls map[string]int
	producerCalls map[string]int

	communication string
}

func defaultNatsConfig() NatsAnalysisConfig {
	return NatsAnalysisConfig{
		consumerCalls: map[string]int{
			"Subscribe": 1,
		},
		// Could not find one.
		producerCalls: map[string]int{
			"NewBatchNotifyMsg":                    1,
			"NewBusinessDeclarationsNotifyMsg":     1,
			"NewCostNotifyMsg":                     1,
			"NewCostReportsNotifyMsg":              1,
			"NewEDSNNotifyMsg":                     1,
			"NewEdsnRequestNotifyMsg":              1,
			"NewEnergyMarketRetrievalNotifyMsg":    1,
			"NewEnodeNotifyMsg":                    1,
			"NewEventReadingNotifyMsg":             1,
			"NewInvoiceNotifyMsg":                  1,
			"NewFeedbackNotifyMsg":                 1,
			"NewBouncedNotifyMsg":                  1,
			"NewClickedNotifyMsg":                  1,
			"NewOpenedNotifyMsg":                   1,
			"NewDeliveredNotifyMsg":                1,
			"NewMailCreatedNotifyMsg":              1,
			"NewMeteringPointEDSNNotifyMsg":        1,
			"NewMeteringPointFailureNotifyMsg":     1,
			"NewPaymentNotifyMsg":                  1,
			"NewPFourIncompleteNotifyMsg":          1,
			"NewPostmarkNotifyMsg":                 1,
			"NewReadingNotifyMsg":                  1,
			"NewSnapshotLossNotifyMsg":             1,
			"NewSnapshotStartdateChangedNotifyMsg": 1,
			"NewSnapshotNotifyMsg":                 1,
			"NewUsersNotifyMsg":                    1,
			"NewPFourNotifyMsg":                    1,
		},
		communication: "NATS",
	}
}
