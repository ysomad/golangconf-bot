package schedule

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Talk struct {
	StartsAt time.Time
	Title    string
	URL      *url.URL
	Speakers []string
	Duration time.Duration
}

const dateLayout = "02/01/2006 15:04:05"

var (
	errNilReader  = errors.New("пустой файл")
	errInvalidCSV = errors.New("невалидный файл")
	errEmptyCol   = errors.New("пустая ячейка")

	errInvalidStartTime = errors.New("некорректное время начала доклада")
	errInvalidDuration  = errors.New("некорректная длительность доклада")
	errInvalidSpeakers  = errors.New("должен быть указан хотя бы один спикер")
	errInvalidURL       = errors.New("некорректный url доклада")

	errTimezoneNotLoaded = errors.New("ошибка загрузки часового пояса")
)

func readCSV(f io.Reader) ([][]string, error) {
	if f == nil {
		return nil, errNilReader
	}

	reader := csv.NewReader(f)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errInvalidCSV, err)
	}

	return records, err
}

// NewFromCSV reads csv file and parses it into schedule of talks,
// returns list of CSVError.
func NewFromCSV(file io.Reader) ([]Talk, []error) {
	records, err := readCSV(file)
	if err != nil {
		return nil, []error{err}
	}

	var (
		offset = 1 // skipped header rows
		talks  = make([]Talk, len(records)-offset)
		errs   []error
	)

	for i, row := range records {
		// skip header row
		if i == 0 {
			continue
		}

		talk := Talk{}

		for j, col := range row {
			if col == "" {
				errs = append(errs, newCSVError(errEmptyCol, i, j))
				continue
			}

			switch j {
			case 0: // talk start time
				loc, err := time.LoadLocation("Europe/Moscow")
				if err != nil {
					errs = append(errs, newCSVError(fmt.Errorf("%w: %w", errTimezoneNotLoaded, err), i, j))
					continue
				}

				talk.StartsAt, err = time.ParseInLocation(dateLayout, col, loc)
				if err != nil {
					errs = append(errs, newCSVError(fmt.Errorf("%w: %w", errInvalidStartTime, err), i, j))
				}
			case 1: // talk duration
				mins, err := strconv.Atoi(col)
				if err != nil {
					errs = append(errs, newCSVError(fmt.Errorf("%w: %w", errInvalidDuration, err), i, j))
					continue
				}

				talk.Duration = time.Minute * time.Duration(mins)
			case 2: // talk title
				talk.Title = col
			case 3: // talk speakers
				rawSpeakers := strings.Split(col, ",")
				if len(rawSpeakers) < 1 {
					errs = append(errs, newCSVError(errInvalidSpeakers, i, j))
					continue
				}

				talk.Speakers = make([]string, len(rawSpeakers))
				for i, s := range rawSpeakers {
					talk.Speakers[i] = strings.TrimSpace(s)
				}
			case 4: // talk url
				col = strings.TrimPrefix(col, "<")
				col = strings.TrimSuffix(col, ">")

				talk.URL, err = url.Parse(col)
				if err != nil {
					errs = append(errs, newCSVError(fmt.Errorf("%w: %w", errInvalidURL, err), i, j))
				}
			}
		}

		talks[i-offset] = talk
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return talks, nil
}
