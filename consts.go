package bugsnagperformance

import "time"

const (
	samplingAttribute          = "bugsnag.sampling.p"
	samplingResponseHeader     = "Bugsnag-Sampling-Probability"
	samplingRequestHeader      = "Bugsnag-Span-Sampling"
	fetcherRetryInterval       = 30 * time.Second
	fetcherRefreshInterval     = 24 * time.Hour
	fetcherRequestBody         = `{"resourceSpans": []}`
	deploymentEnvAttribute     = "deployment.environment"
	serviceNameAttribute       = "service.name"
	serviceVersionAttribute    = "service.version"
	bugsnagSDKNameAttribute    = "bugsnag.telemetry.sdk.name"
	bugsnagSDKVersionAttribute = "bugsnag.telemetry.sdk.version"
	sdkName                    = "Go Bugsnag Performance SDK"
)
