package porcelain

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"text/scanner"
)

var (
	ErrParser          = errors.New("failed to parse")
	ErrUnexpectedEOF   = errors.New("unexpected EOF")
	ErrUnexpectedToken = errors.New("unexpected token")
)

// Parser understands git command output.
type Parser interface {
	ParseVersion() (GitVersion, error)
	ParseBranch() (string, error)
	ParseCommit() (string, error)
	ParseAhead() (int, error)
	ParseBehind() (int, error)
}

// gitParser implements Parser
type gitParser struct {
	rawbuf io.Reader
	scanner.Scanner
	lastError error
}

func NewParser(buf io.Reader) Parser {
	p := gitParser{rawbuf: buf}
	p.Init(buf)
	p.Error = func(_ *scanner.Scanner, msg string) { p.lastError = errors.New(msg) }
	return &p
}

// ParseVersion parser the git command version number and outputs it as an array with 3 uint integers
// Expected input: `git [--version|version]`
//
func (p *gitParser) ParseVersion() (GitVersion, error) {
	var ver = GitVersion([3]int{0, 0, 0})
	var idx = 0

	p.Mode ^= scanner.ScanFloats
	p.Mode |= scanner.ScanInts
	for tok := p.Scan(); tok != scanner.EOF; tok = p.Scan() {
		if tok == scanner.Int {
			i, _ := strconv.Atoi(p.TokenText())
			ver[idx] = i
		}
		if p.TokenText() == "." {
			idx += 1
			continue
		}
	}
	return ver, p.lastError
}

type hdr struct {
	branch, commit string
	ahead, behind  int
}

var cache *hdr

func (p *gitParser) parseHeaders() (*hdr, error) {
	if cache != nil {
		return cache, nil
	}

	cache = new(hdr)
	// p.Whitespace ^= 1 << '\n'
	for tok := p.Scan(); tok != scanner.EOF; tok = p.Scan() {
		// break if not a header
		if p.TokenText() != "#" {
			break
		}
		if err := p.parseHeader(cache); err != nil {
			return nil, err
		}
	}
	return cache, nil
}

func (p *gitParser) parseHeader(header *hdr) error {
	if tok := p.Scan(); tok == scanner.EOF {
		return ErrUnexpectedEOF
	}
	switch p.TokenText() {
	case "branch":
		if err := p.parseBranchHeader(header); err != nil {
			return err
		}
	}
	return nil
}

func (p *gitParser) parseBranchHeader(header *hdr) (err error) {
	if tok := p.Scan(); tok == scanner.EOF {
		return ErrUnexpectedEOF
	}
	if p.TokenText() != "." {
		return fmt.Errorf("%s %q", ErrUnexpectedToken, p.TokenText())
	}
	if tok := p.Scan(); tok == scanner.EOF {
		return ErrUnexpectedEOF
	}
	switch p.TokenText() {
	case "oid":
		header.commit, err = p.parseBranchHeaderField()
		if err != nil {
			return err
		}
	case "head":
		header.branch, err = p.parseBranchHeaderField()
		if err != nil {
			return err
		}
	case "ab":
		header.ahead, header.behind, err = p.parseAheadBehindField()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("%s %q", ErrUnexpectedToken, p.TokenText())
	}
	return nil
}

func (p *gitParser) parseBranchHeaderField() (string, error) {
	if tok := p.Scan(); tok == scanner.EOF {
		return "", ErrUnexpectedEOF
	}
	return p.TokenText(), nil
}

func (p *gitParser) parseAheadBehindField() (ahead int, behind int, err error) {
	for tok := p.Scan(); tok != scanner.EOF; tok = p.Scan() {
		if p.TokenText() == "+" {
			if tok = p.Scan(); tok == scanner.EOF {
				return 0, 0, ErrUnexpectedEOF
			}
			ahead, err = strconv.Atoi(p.TokenText())
			if err != nil {
				return 0, 0, err
			}
		}
		if p.TokenText() == "-" {
			if tok = p.Scan(); tok == scanner.EOF {
				return 0, 0, ErrUnexpectedEOF
			}
			behind, err = strconv.Atoi(p.TokenText())
			if err != nil {
				return 0, 0, err
			}
		}
	}
	return ahead, behind, nil
}

func (p *gitParser) ParseBranch() (string, error) {
	h, err := p.parseHeaders()
	if err != nil {
		return "", err
	}
	return h.branch, nil
}

func (p *gitParser) ParseCommit() (string, error) {
	h, err := p.parseHeaders()
	if err != nil {
		return "", err
	}
	return h.commit, nil
}

func (p *gitParser) ParseAhead() (int, error) {
	h, err := p.parseHeaders()
	if err != nil {
		return 0, err
	}
	return h.ahead, nil
}
func (p *gitParser) ParseBehind() (int, error) {
	h, err := p.parseHeaders()
	if err != nil {
		return 0, err
	}
	return h.behind, nil
}
