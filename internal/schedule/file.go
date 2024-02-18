package schedule

import (
	"encoding/csv"
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

const (
	dateLayout      = "02/01/2006 15:04:05"
	invalidStartsAt = "некорректное время начала доклада, необходимый формат: %s"
	invalidSpeakers = "у доклада должен быть хотя бы один спикер, несколько указываются через ','"
)

func csvError(msg string, row, col int) error {
	return fmt.Errorf("%s (ряд %d, столбец %d)", msg, row+1, col+1)
}

func ParseCSV(file io.Reader) ([]Talk, []error) {
	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, []error{fmt.Errorf("невалидный файл: %w", err)}
	}

	var (
		talks = make([]Talk, len(records)-1)
		errs  []error
	)

	for i, row := range records {
		// if no header or only header present in the file
		if len(row) <= 1 {
			return nil, []error{fmt.Errorf("пустой файл")}
		}

		// skip header row
		if i == 0 {
			continue
		}

		// skip parsing row if there is not 5 columns
		if len(row) != 5 {
			errs = append(errs, csvError("пропущена обработка ряда, должно быть 5 столбцов в каждом ряду", i, 0))
			continue
		}

		talk := Talk{}

		for j, col := range row {
			if col == "" {
				errs = append(errs, csvError("пустая ячейка", i, j))
				continue
			}

			switch j {
			case 0: // talk start time
				loc, err := time.LoadLocation("Europe/Moscow")
				if err != nil {
					errs = append(errs, csvError(err.Error(), i, j))
					continue
				}

				talk.StartsAt, err = time.ParseInLocation(dateLayout, col, loc)
				if err != nil {
					errs = append(errs, csvError(fmt.Sprintf(invalidStartsAt, dateLayout), i, j))
				}
			case 1: // talk duration
				mins, err := strconv.Atoi(col)
				if err != nil {
					errs = append(errs, csvError("длительность доклада должна быть числовой строкой в минутах", i, j))
					continue
				}

				talk.Duration = time.Minute * time.Duration(mins)
			case 2: // talk title
				talk.Title = col
			case 3: // talk speakers
				rawSpeakers := strings.Split(col, ",")
				if len(rawSpeakers) < 1 {
					errs = append(errs, csvError(invalidSpeakers, i, j))
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
					errs = append(errs, csvError("некорректный URL на доклад", i, j))
				}
			}
		}

		talks[i-1] = talk
	}

	fmt.Println(talks)

	return talks, errs
}
