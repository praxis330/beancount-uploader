package main

import (
	"errors"
	"regexp"
	"strings"
)

var (
	DateNotFoundError    = errors.New("Unable to find date in the subject line.")
	SubjectParseError    = errors.New("Subject line must have a date and description.")
	DetailLineParseError = errors.New("Detail line must contain an account, an amount, and a currency.")
)

type Parser interface {
	Parse(r *Request) error
	GetTitle() string
}

func GetBeancountItem(r *Request, p Parser) error {
	return p.Parse(r)
}

type detailLine struct {
	Account  string `json: account`
	Amount   string `json: amount`
	Currency string `json: currency`
}

type BeancountItem struct {
	Date    string       `json: date`
	Subject string       `json: subject`
	Details []detailLine `json: details`
}

func getRemainderString(a []string, index int) string {
	var s string
	if len(a) > index {
		s = strings.Join(a[index:], ` `)
	}
	return s
}

func (b *BeancountItem) getDate(s string) error {
	r, _ := regexp.Compile("2[0-9]{3}-[0-9]{2}-[0-9]{2}")
	match := r.FindString(s)

	if match == "" {
		return DateNotFoundError
	}

	b.Date = match

	return nil
}

func (b *BeancountItem) getSubject(s string) error {
	array := strings.Split(s, ` `)

	if len(array) < 2 {
		return SubjectParseError
	}

	b.Subject = getRemainderString(array, 1)

	return nil
}

func (b *BeancountItem) parseSubjectLine(s string) error {
	error := b.getDate(s)

	if error != nil {
		return error
	}

	error = b.getSubject(s)

	if error != nil {
		return error
	}

	return nil
}

func (b *BeancountItem) parseDetails(s []string) error {

	for _, line := range s {
		err := b.addDetailLine(line)

		if err != nil {
			return err
		}
	}

	return nil
}

func (b *BeancountItem) addDetailLine(l string) error {
	array := strings.Split(l, ` `)

	if len(array) < 3 {
		return DetailLineParseError
	}

	b.Details = append(b.Details, detailLine{array[0], array[1], array[2]})

	return nil
}

func (b *BeancountItem) GetTitle() string {
	return b.Date + `_` + strings.Replace(b.Subject, ` `, `-`, -1)
}

func (b *BeancountItem) Parse(r *Request) error {
	err := r.validateBodyLength()

	if err != nil {
		return err
	}

	subjectString := r.getAndTrimSubjectString()

	err = b.parseSubjectLine(subjectString)

	if err != nil {
		return err
	}

	otherDetails := r.getAndTrimDetails()

	err = b.parseDetails(otherDetails)

	if err != nil {
		return err
	}

	return nil
}
