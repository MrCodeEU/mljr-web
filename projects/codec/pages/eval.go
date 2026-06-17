package pages

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html"
	"net/url"
	"strings"
)

type CodecInput struct {
	Mode  string // hash | base64 | url | html
	Op    string // encode | decode (ignored for hash)
	Input string
}

type HashResult struct {
	MD5    string
	SHA1   string
	SHA256 string
	SHA512 string
}

type CodecResult struct {
	Mode   string
	Op     string
	Input  string
	Output string
	Err    string
	Hash   *HashResult
}

func EvalCodec(inp CodecInput) CodecResult {
	r := CodecResult{Mode: inp.Mode, Op: inp.Op, Input: inp.Input}

	switch inp.Mode {
	case "hash":
		r.Hash = computeHashes(inp.Input)
	case "base64":
		if inp.Op == "decode" {
			b, err := base64.StdEncoding.DecodeString(strings.TrimSpace(inp.Input))
			if err != nil {
				// Try URL-safe variant
				b, err = base64.URLEncoding.DecodeString(strings.TrimSpace(inp.Input))
			}
			if err != nil {
				r.Err = "invalid base64: " + err.Error()
			} else {
				r.Output = string(b)
			}
		} else {
			r.Output = base64.StdEncoding.EncodeToString([]byte(inp.Input))
		}
	case "url":
		if inp.Op == "decode" {
			decoded, err := url.QueryUnescape(inp.Input)
			if err != nil {
				r.Err = "invalid percent-encoding: " + err.Error()
			} else {
				r.Output = decoded
			}
		} else {
			r.Output = url.QueryEscape(inp.Input)
		}
	case "html":
		if inp.Op == "decode" {
			r.Output = html.UnescapeString(inp.Input)
		} else {
			r.Output = html.EscapeString(inp.Input)
		}
	default:
		r.Err = fmt.Sprintf("unknown mode: %q", inp.Mode)
	}

	return r
}

func computeHashes(input string) *HashResult {
	data := []byte(input)
	md5h := md5.Sum(data)
	sha1h := sha1.Sum(data)
	sha256h := sha256.Sum256(data)
	sha512h := sha512.Sum512(data)
	return &HashResult{
		MD5:    hex.EncodeToString(md5h[:]),
		SHA1:   hex.EncodeToString(sha1h[:]),
		SHA256: hex.EncodeToString(sha256h[:]),
		SHA512: hex.EncodeToString(sha512h[:]),
	}
}
