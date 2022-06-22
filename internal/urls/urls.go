package urls

/*[[[cog
import cog
from generateMappings import *
cog.outl(
	generate_golang_url_mapping(
		parse_url_mappings(
			load_raw_url_mappings(),
		),
	),
)
]]]*/
const (
	OK                    = "/ok"
	AddDomainToCrawlQueue = "/job/add"
	CrawlerRequestJob     = "/job/request"
	RequestPreflightCheck = "/page/preflight"
	DigestPageLoad        = "/page/digest"
)

// [[[end]]]
