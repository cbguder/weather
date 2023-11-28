# Weather

What's the "perfect weather" score for your location?

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
Oakland, Alameda County, California, United States
╭─────────────┬──────────────────────────┬───────────┬──────────┬─────────┬───────╮
│ ID          │ Name                     │ Elevation │ Distance │ Records │ Score │
├─────────────┼──────────────────────────┼───────────┼──────────┼─────────┼───────┤
│ USC00046336 │ OAKLAND MUSEUM           │     9.10m │   0.58mi │    3377 │ 75.48 │
│ USC00040693 │ BERKELEY                 │    94.50m │   4.87mi │    3625 │ 69.54 │
│ USR0000COKN │ OAKLAND NORTH CALIFORNIA │   427.60m │   5.02mi │    3718 │ 73.91 │
│ USW00023230 │ OAKLAND INTL AP          │     1.50m │   6.35mi │    3964 │ 67.05 │
│ USR0000COKS │ OAKLAND SOUTH CALIFORNIA │   333.80m │   7.04mi │    3716 │ 73.28 │
│ USW00023272 │ SAN FRANCISCO DWTN       │    45.70m │   8.82mi │    3953 │ 81.94 │
│ USC00047414 │ RICHMOND                 │     6.10m │   9.82mi │    3161 │ 73.30 │
╰─────────────┴──────────────────────────┴───────────┴──────────┴─────────┴───────╯

$ weather trends --after 2020-01-01 USC00046336
╭──────┬───────┬───────┬───────┬───────┬────────┬────────┬────────┬────────┬────────┬───────┬───────┬───────┬───────╮
│ Year │   Jan │   Feb │   Mar │   Apr │    May │    Jun │    Jul │    Aug │    Sep │   Oct │   Nov │   Dec │       │
├──────┼───────┼───────┼───────┼───────┼────────┼────────┼────────┼────────┼────────┼───────┼───────┼───────┼───────┤
│ 2020 │ 51.61 │ 82.76 │ 67.74 │ 86.67 │  87.10 │  90.00 │ 100.00 │  90.32 │  83.33 │ 74.19 │ 70.00 │ 32.26 │ 76.23 │
│ 2021 │ 44.83 │ 45.45 │ 64.52 │ 96.67 │ 100.00 │  93.33 │ 100.00 │  93.55 │  93.33 │ 70.97 │ 83.33 │ 35.48 │ 77.40 │
│ 2022 │ 67.74 │ 67.86 │ 70.97 │ 76.67 │  96.77 │  90.00 │ 100.00 │ 100.00 │  75.86 │ 96.77 │ 53.33 │ 16.13 │ 76.10 │
│ 2023 │ 22.58 │  7.14 │ 22.58 │ 83.33 │  83.87 │ 100.00 │ 100.00 │  90.32 │ 100.00 │ 77.42 │ 80.00 │       │ 69.75 │
├──────┼───────┼───────┼───────┼───────┼────────┼────────┼────────┼────────┼────────┼───────┼───────┼───────┼───────┤
│      │ 46.72 │ 51.40 │ 56.45 │ 85.83 │  91.94 │  93.33 │ 100.00 │  93.55 │  88.24 │ 79.84 │ 70.91 │ 27.96 │ 75.00 │
╰──────┴───────┴───────┴───────┴───────┴────────┴────────┴────────┴────────┴────────┴───────┴───────┴───────┴───────╯
```
