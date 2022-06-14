package natsanalyzer

// NatsAnalysisConfig is a structure holding
// parameters necessary for NATS dependencies
// analysis. Currently, it is rather thin, but it
// was introduced in an effort to make the tool
// more extendable.
type NatsAnalysisConfig struct {
	communication string
}

func defaultNatsConfig() NatsAnalysisConfig {
	return NatsAnalysisConfig{
		communication: "NATS",
	}
}
