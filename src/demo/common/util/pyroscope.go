package util

import (
	"github.com/grafana/pyroscope-go"
	"os"
	"runtime"
)

func InitPyroscope(appName string) {

	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)

	pyroscope.Start(pyroscope.Config{
		ApplicationName: appName,

		// replace this with the address of pyroscope server
		ServerAddress: "http://logtail-kubernetes-metrics.sls-monitoring:4040",

		// you can disable logging by setting this to nil
		Logger: nil,

		// optionally, if authentication is enabled, specify the API key:
		// AuthToken:    os.Getenv("PYROSCOPE_AUTH_TOKEN"),

		// you can provide static tags via a map:
		Tags: map[string]string{"hostname": os.Getenv("HOSTNAME"), "version": "0.0.1", "environment": "test"},

		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})

}
