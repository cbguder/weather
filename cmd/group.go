package cmd

import "github.com/cbguder/weather/noaa"

type recordGroup struct {
	key     string
	records []noaa.DailyRecord
}

func groupRecordsByKey(records []noaa.DailyRecord, keyFunc func(noaa.DailyRecord) string) []recordGroup {
	var orderedKeys []string
	seenKeys := make(map[string]struct{})

	recordsByKey := make(map[string][]noaa.DailyRecord)

	for _, record := range records {
		key := keyFunc(record)

		if _, ok := seenKeys[key]; !ok {
			orderedKeys = append(orderedKeys, key)
			seenKeys[key] = struct{}{}
		}

		recordsByKey[key] = append(recordsByKey[key], record)
	}

	var groups []recordGroup

	for _, key := range orderedKeys {
		groups = append(groups, recordGroup{
			key:     key,
			records: recordsByKey[key],
		})
	}

	return groups
}
