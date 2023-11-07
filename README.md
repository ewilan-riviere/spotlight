# spotlight

[![go][go-version-src]][go-version-href]
[![tests][tests-src]][tests-href]
[![license][license-src]][license-href]

Scan system for system health, powered by Go.

> [!IMPORTANT]
>
> Not stable.

## Install

```bash
go install github.com/ewilan-riviere/spotlight@latest
```

## Usage

### Disk

```bash
spotlight disk
```

Output will be like:

```bash
Filesystem     Type  Size  Used Avail Use% Mounted on
/dev/sda1      ext4  150G   56G   89G  39% /
```

### RAM

```bash
spotlight ram
```

Output will be like:

```bash
RAM: 3.7G / 7.7G
```

### CPU

```bash
spotlight cpu
```

Output will be like:

```bash
CPU: 0.00% / 0.00%
```

### Health

```bash
spotlight health
```

Output will be Disk, RAM and CPU.

### files

```bash
spotlight files
```

Options:

- `--size|s` (default: `30`): minimum size of files to display (in MB)
- `--exts|e` (default: []): extensions to exclude (ex: `--exts=cbz --exts=mp3`)

Output will be like:

```bash
347M	/home/username/.cache/...
211M	/usr/lib/chromium/...
139M	/home/username/.cache/ms-playwright/...
124M	/var/lib/docker/overlay2/...
123M	/var/www/...
119M	/var/lib/docker/overlay2/...
112M	/usr/lib/x86_64-linux-gnu...
110M	/var/lib/docker/overlay2/...
```

### websites

Send a request to websites with `ping` and `curl` to check if they are up.

```bash
spotlight ping -d=example.com -d=example2.com
```

## License

[MIT](LICENSE) © Ewilan Rivière

[go-version-src]: https://img.shields.io/static/v1?style=flat&label=Go&message=v1.21&color=00ADD8&logo=go&logoColor=ffffff&labelColor=18181b
[go-version-href]: https://go.dev/
[tests-src]: https://img.shields.io/github/actions/workflow/status/ewilan-riviere/notifier/run-tests.yml?branch=main&label=tests&style=flat&colorA=18181B
[tests-href]: https://packagist.org/packages/ewilan-riviere/notifier
[license-src]: https://img.shields.io/github/license/ewilan-riviere/spotlight.svg?style=flat&colorA=18181B&colorB=00ADD8
[license-href]: https://github.com/ewilan-riviere/spotlight/blob/main/LICENSE
