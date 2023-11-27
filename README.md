# Weather

What's the "perfect weather" score for your location?

## Usage

```
Usage:
  weather [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  months      Show score by month for a single station
  nearby      Show score for nearby weather stations
  trends      Show historic trends for a single station

Flags:
      --after string    only use records after this date
      --before string   only use records before this date
      --cache string    cache directory (default ~/.cache/weather)
  -h, --help            help for weather

Use "weather [command] --help" for more information about a command.
```

## Example

```
$ weather nearby "Oakland, CA"
Oakland, Alameda County, California, United States
╭─────────────┬──────────────────────────┬───────────┬──────────┬─────────┬───────╮
│ ID          │ Name                     │ Elevation │ Distance │ Records │ Score │
├─────────────┼──────────────────────────┼───────────┼──────────┼─────────┼───────┤
│ USC00046336 │ OAKLAND MUSEUM           │     9.10m │   0.58mi │    3050 │ 53.25 │
│ USC00040693 │ BERKELEY                 │    94.50m │   4.87mi │    3295 │ 40.73 │
│ USR0000COKN │ OAKLAND NORTH CALIFORNIA │   427.60m │   5.02mi │    3388 │ 37.40 │
│ USW00023230 │ OAKLAND INTL AP          │     1.50m │   6.35mi │    3635 │ 46.11 │
│ USR0000COKS │ OAKLAND SOUTH CALIFORNIA │   333.80m │   7.04mi │    3386 │ 37.33 │
│ USW00023272 │ SAN FRANCISCO DWTN       │    45.70m │   8.82mi │    3623 │ 58.54 │
│ USC00047414 │ RICHMOND                 │     6.10m │   9.82mi │    2871 │ 50.40 │
╰─────────────┴──────────────────────────┴───────────┴──────────┴─────────┴───────╯

$ weather months --after 2013-01-01 USC00046336
╭───────────┬─────────┬───────╮
│ Month     │ Records │ Score │
├───────────┼─────────┼───────┤
│ January   │     297 │ 10.77 │
│ February  │     247 │ 20.24 │
│ March     │     279 │ 25.45 │
│ April     │     262 │ 51.15 │
│ May       │     277 │ 81.23 │
│ June      │     271 │ 84.50 │
│ July      │     306 │ 91.83 │
│ August    │     286 │ 84.97 │
│ September │     277 │ 73.29 │
│ October   │     308 │ 63.96 │
│ November  │     290 │ 37.24 │
│ December  │     277 │  9.03 │
╰───────────┴─────────┴───────╯

$ weather trends --after 2020-01-01 USC00046336      
╭──────┬───────┬───────┬───────┬───────┬───────┬────────┬────────┬───────┬───────┬───────┬───────┬───────╮
│ Year │   Jan │   Feb │   Mar │   Apr │   May │    Jun │    Jul │   Aug │   Sep │   Oct │   Nov │   Dec │
├──────┼───────┼───────┼───────┼───────┼───────┼────────┼────────┼───────┼───────┼───────┼───────┼───────┤
│ 2020 │  6.45 │ 20.69 │ 45.16 │ 70.00 │ 67.74 │  76.67 │  90.32 │ 70.97 │ 70.00 │ 54.84 │ 20.00 │ 0.00  │
│ 2021 │ 13.79 │ 13.64 │  6.45 │ 43.33 │ 90.32 │  90.00 │ 100.00 │ 93.55 │ 80.00 │ 48.39 │ 60.00 │ 22.58 │
│ 2022 │  9.68 │ 25.00 │ 16.13 │ 36.67 │ 80.65 │  76.67 │  93.55 │ 90.32 │ 55.17 │ 90.32 │  0.00 │ 0.00  │
│ 2023 │  3.23 │  0.00 │  0.00 │ 26.67 │ 80.65 │ 100.00 │  93.55 │ 77.42 │ 96.67 │ 54.84 │ 45.00 │       │
╰──────┴───────┴───────┴───────┴───────┴───────┴────────┴────────┴───────┴───────┴───────┴───────┴───────╯
```
