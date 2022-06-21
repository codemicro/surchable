"""[[[cog
import cog
from generateMappings import *
cog.outl(
	generate_python_url_mapping(
		parse_url_mappings(
			load_raw_url_mappings(),
		),
	),
)
]]]"""
OK = "/ok"
ADD_DOMAIN_TO_CRAWL_QUEUE = "/job/add"
CRAWLER_REQUEST_JOB = "/job/request"
REQUEST_PREFLIGHT_CHECK = "/page/preflight"

# [[[end]]]
