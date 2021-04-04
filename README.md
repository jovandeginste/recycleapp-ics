Generate an ical file from the recycleapp.be API

```bash
./recycleapp-ics \
  -zipcode "your zip code" \
  -street "your street name" \
  -house "your house number" > cal.ics
```

Example:

```bash
./recycleapp-ics -zipcode 1000 -street "Nieuwstraat" -house 1 > cal.ics
```

Optional parameter:

- `lang`: pick your language
  - options: nl, fr, en, de
  - default: nl
