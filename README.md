# Weather

What's the "perfect weather" score for your location?

## Scoring

The "perfect weather" score for a given time period is calculated as the
percentage of "good days" for which records are available within the time period.

When using daily records (the default), a "good day" is defined as a day with:

* Temperatures between 45°F and 85°F
* 1mm of precipitation or less

When using hourly records, a "good day" is defined as a day with 6 or more
consecutive "good hours," where a "good hour" is defined as an hour with:

* Temperatures between 50°F and 85°F
* 0.2mm of precipitation or less

## Usage

```
Usage:
  weather [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  nearby      Show score for nearby weather stations
  trends      Show historic trends for a single station

Flags:
      --after string    only use records after this date
      --before string   only use records before this date
      --cache string    cache directory (default ~/.cache/weather)
  -h, --help            help for weather
      --ignore-cache    ignore cached data

Use "weather [command] --help" for more information about a command.
```

## Example

```
$ weather nearby --after 2013-01-01 "Oakland, CA"
╭─────────────────────────────────────────────────────────────────────────────────╮
│ Stations near Oakland, Alameda County, California, United States                │
├─────────────┬──────────────────────────┬───────────┬──────────┬─────────┬───────┤
│ ID          │ Name                     │ Elevation │ Distance │ Records │ Score │
├─────────────┼──────────────────────────┼───────────┼──────────┼─────────┼───────┤
│ USC00046336 │ OAKLAND MUSEUM           │     9.10m │   0.58mi │    3410 │ 75.43 │
│ USC00040693 │ BERKELEY                 │    94.50m │   4.87mi │    3665 │ 69.69 │
│ USR0000COKN │ OAKLAND NORTH CALIFORNIA │   427.60m │   5.02mi │    3744 │ 73.90 │
│ USW00023230 │ OAKLAND INTL AP          │     1.50m │   6.35mi │    3997 │ 66.70 │
│ USR0000COKS │ OAKLAND SOUTH CALIFORNIA │   333.80m │   7.04mi │    3742 │ 73.22 │
│ USW00023272 │ SAN FRANCISCO DWTN       │    45.70m │   8.82mi │    3987 │ 81.89 │
│ USC00047414 │ RICHMOND                 │     6.10m │   9.82mi │    3191 │ 72.99 │
╰─────────────┴──────────────────────────┴───────────┴──────────┴─────────┴───────╯

$ weather trends --after 2020-01-01 USC00046336
╭───────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│ Historic Trends for OAKLAND MUSEUM                                                                                │
├──────┬───────┬───────┬───────┬───────┬────────┬────────┬────────┬────────┬────────┬───────┬───────┬───────┬───────┤
│ Year │   Jan │   Feb │   Mar │   Apr │    May │    Jun │    Jul │    Aug │    Sep │   Oct │   Nov │   Dec │       │
├──────┼───────┼───────┼───────┼───────┼────────┼────────┼────────┼────────┼────────┼───────┼───────┼───────┼───────┤
│ 2020 │ 51.61 │ 82.76 │ 67.74 │ 86.67 │  87.10 │  90.00 │ 100.00 │  90.32 │  83.33 │ 74.19 │ 70.00 │ 32.26 │ 76.23 │
│ 2021 │ 44.83 │ 45.45 │ 64.52 │ 96.67 │ 100.00 │  93.33 │ 100.00 │  93.55 │  93.33 │ 70.97 │ 83.33 │ 35.48 │ 77.40 │
│ 2022 │ 67.74 │ 67.86 │ 70.97 │ 76.67 │  96.77 │  90.00 │ 100.00 │ 100.00 │  75.86 │ 96.77 │ 53.33 │ 16.13 │ 76.10 │
│ 2023 │ 22.58 │  7.14 │ 22.58 │ 83.33 │  83.87 │ 100.00 │ 100.00 │  90.32 │ 100.00 │ 77.42 │ 76.67 │ 69.57 │ 69.75 │
├──────┼───────┼───────┼───────┼───────┼────────┼────────┼────────┼────────┼────────┼───────┼───────┼───────┼───────┤
│      │ 46.72 │ 51.40 │ 56.45 │ 85.83 │  91.94 │  93.33 │ 100.00 │  93.55 │  88.24 │ 79.84 │ 70.83 │ 36.21 │ 74.88 │
╰──────┴───────┴───────┴───────┴───────┴────────┴────────┴────────┴────────┴────────┴───────┴───────┴───────┴───────╯
```

For hourly data, retrieve station ID manually from https://www.ncei.noaa.gov/access/search/data-search/global-hourly

```
$ weather trends --after 2020-01-01 --hourly 99849699999
╭─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│ Historic Trends for 99849699999                                                                                     │
├──────┬───────┬───────┬───────┬───────┬────────┬────────┬────────┬────────┬────────┬────────┬────────┬───────┬───────┤
│ Year │   Jan │   Feb │   Mar │   Apr │    May │    Jun │    Jul │    Aug │    Sep │    Oct │    Nov │   Dec │       │
├──────┼───────┼───────┼───────┼───────┼────────┼────────┼────────┼────────┼────────┼────────┼────────┼───────┼───────┤
│ 2020 │ 48.39 │ 72.41 │ 64.52 │ 86.67 │ 100.00 │ 100.00 │ 100.00 │  96.77 │ 100.00 │ 100.00 │  60.00 │ 32.26 │ 79.89 │
│ 2021 │ 43.33 │ 42.31 │ 48.00 │ 70.00 │ 100.00 │ 100.00 │ 100.00 │  96.77 │ 100.00 │ 100.00 │ 100.00 │ 41.94 │ 79.49 │
│ 2022 │ 35.48 │ 46.43 │ 74.19 │ 83.33 │  90.32 │ 100.00 │  93.55 │ 100.00 │ 100.00 │  61.29 │  33.33 │ 16.13 │ 69.59 │
│ 2023 │ 48.39 │ 14.29 │ 29.03 │ 60.00 │ 100.00 │ 100.00 │ 100.00 │ 100.00 │ 100.00 │ 100.00 │  93.33 │ 58.33 │ 75.98 │
├──────┼───────┼───────┼───────┼───────┼────────┼────────┼────────┼────────┼────────┼────────┼────────┼───────┼───────┤
│      │ 43.90 │ 44.14 │ 54.24 │ 75.00 │  97.58 │ 100.00 │  98.35 │  98.39 │ 100.00 │  90.32 │  71.67 │ 35.90 │ 76.21 │
╰──────┴───────┴───────┴───────┴───────┴────────┴────────┴────────┴────────┴────────┴────────┴────────┴───────┴───────╯
```
