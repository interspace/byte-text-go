package validate

import "regexp"

//
//    # These URL validation pattern strings are based on the ABNF from RFC 3986
const (
	validateURLUnreserved = `[a-z0-9\-._~]`
	validateURLPctEncoded = `(?:%[0-9-a-fA-F]{2})`
	validateURLSubDelims  = `[!$&'()*+,;=]`
	validateURLPchar      = `(?:` +
		validateURLUnreserved + `|` +
		validateURLPctEncoded + `|` +
		validateURLSubDelims + `|` +
		`[:\|@])`

	validateURLUserinfo = `(?:` +
		validateURLUnreserved + `|` +
		validateURLPctEncoded + `|` +
		validateURLSubDelims + `|` +
		`:)*`

	validateURLDecOctet = `(?:[0-9]|(?:[1-9][0-9])|(?:1[0-9]{2})|(?:2[0-4][0-9])|(?:25[0-5]))`

	validateURLIpv4 = `(?:` +
		validateURLDecOctet +
		`(?:\.` + validateURLDecOctet + `){3}` +
		`)`

	// Punting on real IPv6 validation for now
	validateURLIpv6 = `(?:\[[a-fA-F0-9:\.]+\])`

	// Also punting on IPvFuture for now
	validateURLIp = `(?:` +
		validateURLIpv4 + `|` + validateURLIpv6 +
		`)`

	// This is more strict than the rfc specifies
	validateURLSubdomainSegment = `(?:[a-z0-9](?:[a-z0-9_\-]*[a-z0-9])?)`
	validateURLDomainSegment    = `(?:[a-z0-9](?:[a-z0-9\-]*[a-z0-9])?)`
	validateURLDomainTld        = `(?:[a-z](?:[a-z0-9\-]*[a-z0-9])?)`
	validateURLDomain           = `(?:(?:` +
		validateURLSubdomainSegment + `\.)*` +
		`(?:` + validateURLDomainSegment + `\.)` +
		validateURLDomainTld + `)`

	validateURLHost = `(?:` + validateURLIp + `|` + validateURLDomain + `)`

	// Unencoded internationalized domains - this doesn't check for invalid UTF-8 sequences
	validateURLUnicodeSubdomainSegment = `(?:` +
		`(?:[a-z0-9]|[^\x00-\x7f])(?:(?:[a-z0-9_\-]` +
		`|[^\x00-\x7f])*(?:[a-z0-9]|[^\x00-\x7f]))?)`

	validateURLUnicodeDomainSegment = `(?:` +
		`(?:[a-z0-9]|[^\x00-\x7f])(?:(?:[a-z0-9\-]|` +
		`[^\x00-\x7f])*(?:[a-z0-9]|[^\x00-\x7f]))?)`

	validateURLUnicodeDomainTld = `(?:` +
		`(?:[a-z]|[^\x00-\x7f])(?:(?:[a-z0-9\-]|` +
		`[^\x00-\x7f])*(?:[a-z0-9]|[^\x00-\x7f]))?)`

	validateURLUnicodeDomain = `(?:` +
		`(?:` + validateURLUnicodeSubdomainSegment + `\.)*` +
		`(?:` + validateURLUnicodeDomainSegment + `\.)` +
		validateURLUnicodeDomainTld + `)`

	validateURLUnicodeHost = `(?:` +
		validateURLIp + `|` +
		validateURLUnicodeDomain +
		`)`

	validateURLPort = `[0-9]{1,5}`

	validateURLUnicodeAuthority = `\A(?:` +
		`(` + validateURLUserinfo + `)@)?` + // $1 userinfo
		`(` + validateURLUnicodeHost + `)` + // $2 host
		`(?::(` + validateURLPort + `))?\z` // $3 port

	validateURLAuthority = `\"(?:` +
		`(` + validateURLUserinfo + `)@)?` + // $1 userinfo
		`(` + validateURLHost + `)` + // $2 host
		`(?::(` + validateURLPort + `))?\z` // $3 port

	validateURLScheme   = `\A(?:[a-z][a-z0-9+\-.]*)\z`
	validateURLPath     = `\A(/` + validateURLPchar + `*)*\z`
	validateURLQuery    = `\A(` + validateURLPchar + `|/|\?)*\z`
	validateURLFragment = `\A(` + validateURLPchar + `|/|\?)*\z`

	// Modified version of RFC 3986 Appendix B
	validateURLUnencoded = `\A` + // Full URL
		`(?:` +
		`([^:/?#]+)://` + // $1 Scheme
		`)?` +
		`([^/?#]*)` + // $2 Authority
		`([^?#]*)` + // $3 Path
		`(?:` +
		`\?([^#]*)` + // $4 Query
		`)?` +
		`(?:` +
		`\#(.*)` + // $5 Fragment
		`)?\z`

	validateURLUnencodedGroupScheme    = 1
	validateURLUnencodedGroupAuthority = 2
	validateURLUnencodedGroupPath      = 3
	validateURLUnencodedGroupQuery     = 4
	validateURLUnencodedGroupFragment  = 5
)

var (
	validateURLUnencodedRe        = regexp.MustCompile(`(?i)` + validateURLUnencoded)
	validateURLSchemeRe           = regexp.MustCompile(`(?i)` + validateURLScheme)
	validateURLPathRe             = regexp.MustCompile(`(?i)` + validateURLPath)
	validateURLQueryRe            = regexp.MustCompile(`(?i)` + validateURLQuery)
	validateURLFragmentRe         = regexp.MustCompile(`(?i)` + validateURLFragment)
	validateURLAuthorityRe        = regexp.MustCompile(`(?i)` + validateURLAuthority)
	validateURLUnicodeAuthorityRe = regexp.MustCompile(`(?i)` + validateURLUnicodeAuthority)
	protocolRe                    = regexp.MustCompile(`(?i)\Ahttps?\z`)
)
