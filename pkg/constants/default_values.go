package constants

const (
	RUNTIME_ROOT         = "/home/user/src/runtime/data"
	RUNTIME_SUBDIRECTORY = "ratelimit"

	DEFAULT_RATELIMITER_SIZE int32 = 1
	DEFAULT_RATELIMITER_PORT int32 = 8081

	REDIS_PORT int32 = 6379

	DEFAULT_STATSD_PORT         = 9125
	DEFAULT_STATSD_MAPPING_DIR  = "/tmp"
	DEFAULT_STATSD_MAPPING_FILE = "statsd_mapping.yml"
	DEFALT_STATSD_LOGLEVEL      = "info"
)
