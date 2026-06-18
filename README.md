# Recycleapp iCalendar Generator

Generate an iCalendar (ICS) or JSON file from the recycleapp.be API for your garbage collection schedule.

## Build Instructions

To compile the binary, run:

```bash
go build ./cmd/recycleapp-ics
```

This will produce the `recycleapp-ics` executable in your current directory.

## Usage

Run the compiled binary with your details:

```bash
./recycleapp-ics \
  --zipcode "your zip code" \
  --street "your street name" \
  --house "your house number" > cal.ics
```

### Example:

```bash
./recycleapp-ics --zipcode 1000 --street "Nieuwstraat" --house 1 > cal.ics
```

### Optional parameters:

- `--lang`: pick your language
  - options: nl, fr, en, de
  - default: nl
- `--format`: output format
  - options: json, ics
  - default: json
- `--year`: the year for the calendar
  - default: current year
