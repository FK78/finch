<p align="center">
  <img src="logo.svg" alt="finch" width="400"/>
</p>

<p align="center">
  Track your portfolio without leaving the terminal. Live prices and P&L, powered by Go.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/go-1.21+-00D4AA?style=flat-square&logo=go&logoColor=white"/>
  <img src="https://img.shields.io/badge/license-MIT-00D4AA?style=flat-square"/>
  <img src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux-00D4AA?style=flat-square"/>
</p>

---

## What is finch?

finch is a no-frills portfolio tracker for the terminal. Tell it what you own and what you paid - it handles the rest. Live prices, real P&L, no browser required.

Built for people who live in the terminal and don't want to context-switch to a browser just to check whether they're up or down.

---

## Features

- Add and manage holdings locally
- Fetch live prices via [Finnhub](https://finnhub.io)
- View P&L per holding and overall
- Allocation breakdown across your portfolio
- Zero accounts, zero dashboards, zero noise

---

## Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- A free [Finnhub API key](https://finnhub.io/register)

---

## Installation

```bash
git clone https://github.com/yourusername/finch.git
cd finch
go build -o finch
```

Move the binary somewhere on your `$PATH` to use it from anywhere:

```bash
mv finch /usr/local/bin/
```

---

## Setup

On first run, finch will prompt you for your Finnhub API key. This is stored locally at `~/.finch/config.json` and never leaves your machine.

To update your API key at any time:

```bash
finch api-key
```

---

## Usage

```bash
finch add       # Add a new holding
finch list      # View all holdings with live prices and P&L
finch remove    # Remove a holding
```

---

## How it works

finch stores your holdings in `~/.finch/holdings.json`. When you run `finch list`, it fetches the current price for each ticker from Finnhub and calculates your P&L based on your buy price and quantity.

Your data never leaves your machine except for the price requests to Finnhub.

---

## License

MIT