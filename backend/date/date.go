package date

import "fmt"

// `Date` represents a date, accurate to the day of a month of a year.
type Date struct {
	Day   uint `json:"day"`
	Month uint `json:"month"`
	Year  uint `json:"year"`
}

// Months of the year.
const (
	Jan = iota + 1
	Feb
	Mar
	Apr
	May
	Jun
	Jul
	Aug
	Sep
	Oct
	Nov
	Dec

	InvalidMonth = "INVALID_MONTH"
)

// MonthToStr converts a month into its string representation. The returned
// string is the full month name. To get a short version of the name, take a
// substring of the first 3 characters.
func MonthToStr(m uint) string {
	switch m {
	case Jan:
		return "January"
	case Feb:
		return "February"
	case Mar:
		return "March"
	case Apr:
		return "April"
	case Jun:
		return "June"
	case Jul:
		return "July"
	case Aug:
		return "August"
	case Sep:
		return "September"
	case Oct:
		return "October"
	case Nov:
		return "November"
	case Dec:
		return "December"
	default:
		return InvalidMonth
	}
}

// `New` is used to compose a new Date. Always use this method to create date
// instances so that all dates in the system are valid.
func New(year, month, day uint) (*Date, error) {
	if !(1 <= month && month <= 12) {
	}
	d := &Date{
		Day:   day,
		Month: month,
		Year:  year,
	}
	if err := d.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid date: %s", err.Error())
	}
	return d, nil
}

// `isLeapYear` returns whether the given year is a leap year using the
// algorithm based on divisibility rules.
func isLeapYear(y uint) bool {
	if y%400 == 0 {
		return true
	} else if y%100 == 0 {
		return false
	} else if y%4 == 0 {
		return true
	} else {
		return false
	}
}

// `IsValid` returns whether the `Date`, `d`, is a valid date. It checks that
// the month and day values are correct as well as whether the day value is
// correct, based on the month and whether the year is a leap one.
// If the returned error is `nil`, then the date is valid, otherwise the error
// contains the reason as to why the date is not valid.
func (d *Date) IsValid() error {
	// sanity checks on the day and month values
	if d.Day == 0 || d.Day > 31 {
		return fmt.Errorf(
			"expected day (%d) to be between 0 and 31 (inclusive)",
			d.Day,
		)
	} else if d.Month == 0 || d.Month > 12 {
		return fmt.Errorf(
			"expected month (%d) to be between 0 and 12 (inclusive)",
			d.Month,
		)
	}
	// checks to ensure that the day value is valid, based on the month
	isLeap := isLeapYear(d.Year)
	var ub uint
	switch d.Month {
	case Feb:
		{
			if isLeap {
				if d.Day > 29 {
					return fmt.Errorf(
						"day (%d) greater than 29 in leap year (%d)",
						d.Day, d.Year,
					)
				}
			} else {
				if d.Day > 28 {
					return fmt.Errorf(
						"day (%d) greater than 28 in non-leap year (%d)",
						d.Day, d.Year,
					)
				}
			}
		}
	case Jan, Mar, May, Jul, Aug, Oct, Dec:
		ub = 31
	default:
		ub = 30
	}
	if d.Day > ub {
		return fmt.Errorf(
			"expected day (%d) to be at most %d for month %s",
			d.Day, ub, MonthToStr(d.Month),
		)
	}
	return nil
}

// `String` returns a string representation of the date. For example, the
// `String` method of the date returned by New(1999, 12, 28) would return the
// representation "28th December, 1999". The ordinal representation of the day
// is also taken into consideration. So for the date returned by New(1999, 1,
// 1), the string representation would be "1st January, 1999".
func (d *Date) String() string {
	var ordinalExt string
	switch d.Day % 10 {
	case 1:
		ordinalExt = "st"
	case 2:
		ordinalExt = "nd"
	case 3:
		ordinalExt = "rd"
	default:
		ordinalExt = "th"
	}
	return fmt.Sprintf(
		"%d%s %s, %d",
		d.Day, ordinalExt, MonthToStr(d.Month), d.Year,
	)
}
