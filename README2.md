# Baroo Investor

A small Go desktop app (built with [Fyne](https://fyne.io)) that shows Baro
Ki'Teer's current stock alongside 10‑day platinum price/volume averages from
warframe.market, plus a platinum‑per‑ducat ratio.

## Layout

- **Left sidebar** — a single nav item, "Baroo Investor" with a coin icon.
  Click it to load/refresh the data. (This is intentionally the only thing
  in the sidebar for now, per the spec — easy to add more entries later in
  `ui.go`'s `buildSidebar`.)
- **Right panel** — a status/header line, a progress bar while fetching, and
  a table with columns: Item, Ducats, Avg Platinum (10d), Avg Volume (10d),
  Plat/Ducat, Data points.

## Data sources

1. **Baro's current inventory** — `https://api.warframestat.us/pc/voidTrader`
   (WFCD's community worldstate API). Gives each item's name, ducat cost,
   and credit cost.
2. **warframe.market catalogue** — `GET /v1/items` to map each item's display
   name to its `url_name` slug.
3. **warframe.market statistics** — `GET /v1/items/{url_name}/statistics`.
   This is the legacy v1 endpoint; warframe.market has a newer v2 API in
   progress but as of this writing v2 does not yet expose historical daily
   statistics, so v1 (deprecated but functional) is used here. If it gets
   shut off, swap the implementation in `marketapi.go`.

## Reference window

The app computes averages over the 10 days ending on Baro's arrival date
(inclusive). By default that arrival date comes straight from the
`voidTrader.activation` field returned by the worldstate API.

For reproducibility, `main.go` currently **pins** the reference date to
`2026-07-10` via the `overrideStockDate` variable, which gives the
1.07.2026–10.07.2026 window described in the spec. Set that variable to
`time.Time{}` to always use whatever Baro is *actually* selling right now
instead of the fixed test date.

## Rate limiting

warframe.market asks clients to stay around 3 requests/second. The app
pauses ~350ms between statistics calls to stay comfortably under that,
fetching items sequentially. With Baro's usual ~40-50 item stock, a full
refresh takes roughly 15-20 seconds.

## Icon

The sidebar coin is an original stylized SVG (a gold hexagon with a "D"),
not a copy of the in-game ducat icon — swap `ducatIconSVG` in `ui.go` for
your own asset if you have one you're licensed to use.

## Building & running

You need Go 1.21+ and a working C toolchain (Fyne uses cgo for the GL
backend on desktop). This app uses `fyne.Do`, which needs Fyne v2.6+ (pinned
in `go.mod`):

```bash
cd baro-investor
go mod tidy      # downloads fyne.io/fyne/v2 and its deps
go run .
```

To build a standalone binary:

```bash
go build -o baro-investor .
```

On Linux you may need `gcc`, `libgl1-mesa-dev`, `xorg-dev` (or your distro's
equivalents) installed for the OpenGL/X11 bindings Fyne uses. On Windows,
a working `gcc` (e.g. via MSYS2/TDM-GCC) is required for cgo.

## Known limitations / next steps

- Item name matching between the two APIs is done by normalizing to
  lowercase alphanumerics. Baro's stock is almost entirely Prime
  parts/sets/blueprints, which match cleanly, but if a mismatch ever occurs
  the row will show "no data" for platinum/volume/ratio while still showing
  ducat cost.
- No caching/persistence yet — every click re-fetches everything.
- No sorting/filtering on the table yet (e.g. sort by Plat/Ducat to find the
  best flips) — the `widget.Table` in `ui.go` is a natural place to add that.
