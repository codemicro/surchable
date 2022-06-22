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

// The below was generated. Do not edit.
// Modify mappings/urls instead.

const (
	OK                    = "/ok"
	AddDomainToCrawlQueue = "/job/add"
	CrawlerRequestJob     = "/job/request"
	RequestPreflightCheck = "/page/preflight"
	DigestPageLoad        = "/page/digest"
)

// [[[end]]]
