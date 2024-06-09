package internal

const (
	ContentTypeAppJSON   = "application/json"
	ContentTypeAppXML    = "application/xml"
	ContentTypeTextPlain = "text/plain"
)

const (
	HeaderAuthorization          = "Authorization"
	HeaderAzureAsync             = "Azure-AsyncOperation"
	HeaderContentLength          = "Content-Length"
	HeaderContentType            = "Content-Type"
	HeaderFakePollerStatus       = "Fake-Poller-Status"
	HeaderLocation               = "Location"
	HeaderOperationLocation      = "Operation-Location"
	HeaderRetryAfter             = "Retry-After"
	HeaderRetryAfterMS           = "Retry-After-Ms"
	HeaderUserAgent              = "User-Agent"
	HeaderWWWAuthenticate        = "WWW-Authenticate"
	HeaderAuxiliaryAuthorization = "x-ms-authorization-auxiliary"
	HeaderXMSClientRequestID     = "x-ms-client-request-id"
	HeaderXMSRequestID           = "x-ms-request-id"
	HeaderXMSErrorCode           = "x-ms-error-code"
	HeaderXMSRetryAfterMS        = "x-ms-retry-after-ms"
)

const BearerTokenPrefix = "Bearer "
