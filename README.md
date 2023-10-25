# spotlight

[![go][go-version-src]][go-version-href]
[![license][license-src]][license-href]

Scan system like disk usage or big files.

## Install

```bash
go install github.com/ewilan-riviere/spotlight@latest
```

## Usage

### Disk usage

```bash
spotlight disk-usage
```

Output will be like:

```bash
Filesystem     Type  Size  Used Avail Use% Mounted on
/dev/sda1      ext4  150G   56G   89G  39% /
```

### Big files

```bash
spotlight big-files
```

Options:

- `--size|s` (default: `30`): minimum size of files to display (in MB)
- `--exts|e` (default: []): extensions to exclude (ex: `--exts=cbz --exts=mp3`)

Output will be like:

```bash
$ sudo find / -xdev -type f -size +30M -exec du -sh {} ';' | sort -rh | head -n50
347M	/home/username/.cache/...
211M	/usr/lib/chromium/...
139M	/home/username/.cache/ms-playwright/...
124M	/var/lib/docker/overlay2/...
123M	/var/www/...
119M	/var/lib/docker/overlay2/...
112M	/usr/lib/x86_64-linux-gnu...
110M	/var/lib/docker/overlay2/...
```

## License

[MIT](LICENSE) © Ewilan Rivière

[go-version-src]: https://img.shields.io/static/v1?style=flat&label=Go&message=v1.21&color=00ADD8&logo=go&logoColor=ffffff&labelColor=18181b
[go-version-href]: https://go.dev/
[license-src]: https://img.shields.io/github/license/ewilan-riviere/spotlight.svg?style=flat&colorA=18181B&colorB=00ADD8
[license-href]: https://github.com/ewilan-riviere/spotlight/blob/main/LICENSE
