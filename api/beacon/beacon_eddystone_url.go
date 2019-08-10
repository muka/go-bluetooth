package beacon

// Taken from https://github.com/suapapa/go_eddystone
// If PR is accepted, move to exposed API
// https://github.com/suapapa/go_eddystone/pull/2

import "errors"

var urlSchemePrefix = []string{
	"http://www.",
	"https://www.",
	"http://",
	"https://",
}

var urlEncoding = []string{
	".com/",
	".org/",
	".edu/",
	".net/",
	".info/",
	".biz/",
	".gov/",
	".com",
	".org",
	".edu",
	".net",
	".info",
	".biz",
	".gov",
}

func decodeURL(prefix byte, encodedURL []byte) (string, error) {
	if int(prefix) >= len(urlSchemePrefix) {
		return "", errors.New("invaild prefix")
	}

	s := urlSchemePrefix[prefix]

	for _, b := range encodedURL {
		switch {
		case 0x00 <= b && b <= 0x13:
			s += urlEncoding[b]
		case 0x0e <= b && b <= 0x20:
			fallthrough
		case 0x7f <= b && b <= 0xff:
			return "", errors.New("invalid byte")
		default:
			s += string(b)
		}
	}

	return s, nil
}
