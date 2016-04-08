package main

import (
	"errors"
	"strings"
)

var (
	BodyParseError = errors.New("Body must have at least three lines separated by a semi-colon.")
)

type Request struct {
	Body   string `json:"body-plain" form:"body-plain" binding:"required"`
	Sender string `json:"sender" form:"sender" binding:"required"`
}

func (r *Request) getAndTrimSubjectString() string {
	subjectLine := r.getLines()[0]

	return strings.Trim(subjectLine, ` `)
}

func (r *Request) getAndTrimDetails() []string {
	otherStrings := r.getLines()[1:]

	for index, str := range otherStrings {
		otherStrings[index] = strings.Trim(str, ` `)
	}

	return otherStrings
}

func (r *Request) validateBodyLength() error {
	lines := r.getLines()

	if len(lines) < 3 {
		return BodyParseError
	}

	return nil
}

func (r *Request) getLines() []string {
	return strings.Split(r.Body, `;`)
}
