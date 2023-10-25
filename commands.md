## CPU

Linux CPU usage

```bash
cat /proc/loadavg | awk '{print $1*100 "%"}'
```

macOS CPU usage

```bash
top -l 2 | grep -E "^CPU"
```

Windows CPU usage

```bash
wmic cpu get loadpercentage | awk 'NR==2 {print $1 "%"}'
```

## Memory

Linux memory usage

```bash
echo "$(free -m | awk 'NR==2 { printf "%.2f GB\n", $3/1024 }') / $(free -m | awk 'NR==2 { printf "%.2f GB\n", $2/1024 }')"
```

macOS memory usage

```bash
top -l 1 -s 0 | awk '/PhysMem/'
```

Windows memory usage

```bash
wmic OS get FreePhysicalMemory,TotalVisibleMemorySize /Value | awk -F"=" 'NR==1 { printf "%.2f GB\n", $2/1024/1024 } NR==2 { printf "%.2f GB\n", $2/1024/1024 }'
```
