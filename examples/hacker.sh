#!/bin/bash

# Fake hacker script — purely cosmetic, does nothing real
# Uses only `gloss` for all formatting (no figlet, toilet, lolcat, pv)

GLOSS="/Users/noraalalou/projects/gloss/gloss"

# macOS-compatible random helpers (no shuf needed)
randint() { echo $(( RANDOM % ($2 - $1 + 1) + $1 )); }
randip() { echo "$(randint 1 255).$(randint 1 255).$(randint 1 255).$(randint 1 255)"; }
randhex() { cat /dev/urandom | LC_ALL=C tr -dc 'a-f0-9' | head -c ${1:-64}; }
randword() {
  local words=("PHANTOM" "GHOST" "SHADOW" "CIPHER" "NEXUS" "VORTEX" "SPECTRE" "DAEMON" "WRAITH" "CHIMERA")
  echo "${words[$((RANDOM % ${#words[@]}))]}"
}

G='\033[0;32m'
LG='\033[1;32m'
R='\033[0;31m'
LR='\033[1;31m'
Y='\033[1;33m'
C='\033[0;36m'
D='\033[2m'
NC='\033[0m'

slow() {
  local text="$1"
  for (( i=0; i<${#text}; i++ )); do
    printf '%s' "${text:$i:1}"
    sleep 0.008
  done
  echo
}

spinner() {
  local msg="$1" duration="$2"
  local frames=('⠋' '⠙' '⠹' '⠸' '⠼' '⠴' '⠦' '⠧' '⠇' '⠏')
  local end=$((SECONDS + duration))
  while [ $SECONDS -lt $end ]; do
    for f in "${frames[@]}"; do
      printf "\r  ${C}$f ${D}$msg${NC}"
      sleep 0.08
    done
  done
  printf "\r\033[2K"
  $GLOSS badge "$msg" --type=success
}

clear

# ── BOOT SEQUENCE ────────────────────────────────────────────────
$GLOSS "SYSTEM BOOT" --font=shadow --gradient=matrix --border=none
echo ""
slow "$(printf "${G}[  OK  ]${NC} Booting secure kernel v6.9.0-darknet-amd64...")"
sleep 0.1
slow "$(printf "${G}[  OK  ]${NC} Loading encrypted overlay filesystem (AES-256-XTS)...")"
sleep 0.1
slow "$(printf "${G}[  OK  ]${NC} Mounting /dev/nvme0n1p2 at /mnt/shadow...")"
sleep 0.1
slow "$(printf "${G}[  OK  ]${NC} Initializing Tor onion routing stack...")"
sleep 0.1
slow "$(printf "${G}[  OK  ]${NC} VPN tunnel established → $(randip) [NordVPN DE#847]")"
sleep 0.1
slow "$(printf "${Y}[ WARN ]${NC} Intrusion detection module active — stealth mode ON")"
sleep 0.1
slow "$(printf "${Y}[ WARN ]${NC} MAC address spoofed: $(randhex 12 | sed 's/../&:/g;s/:$//')")"
sleep 0.2
slow "$(printf "${G}[  OK  ]${NC} Identity scrubbed. You don't exist.")"
sleep 0.3

echo ""
$GLOSS divider --style=heavy --gradient=matrix

# ── RANDOM HEX STREAM ────────────────────────────────────────────
$GLOSS "DECRYPTING" --font=doom --gradient=rainbow --border=none
printf "${D}"
for i in $(seq 1 20); do
  printf "  %s  %s\n" "$(randhex 32)" "$(randhex 32)"
  sleep 0.03
done
printf "${NC}\n"

# ── FAKE PASSWORD CRACKER ─────────────────────────────────────────
$GLOSS divider "PASSWORD CRACKER" --gradient=fire
echo ""
$GLOSS badge "Loading hash dictionary (rockyou.txt — 14,344,391 entries)" --type=warning
sleep 0.3
printf "${D}\n"
wordlist=("password123" "iloveyou" "sunshine" "monkey" "dragon")
hashes=(
  "5f4dcc3b5aa765d61d8327deb882cf99"
  "e10adc3949ba59abbe56e057f20f883e"
  "25d55ad283aa400af464c76d713c07ad"
  "fcea920f7412b5da7be0cf42b8c93759"
  "8621ffdbc5698829397d97767ac13db3"
)
for i in "${!hashes[@]}"; do
  printf "  Cracking: ${Y}%s${NC} " "${hashes[$i]}"
  sleep 0.4
  printf "→ ${LG}%-14s${NC} [MD5 CRACKED]\n" "${wordlist[$i]}"
  sleep 0.15
done
printf "${NC}\n"
$GLOSS bar 100 --style=blocks --label="Cracked" --width=60 --gradient=fire
$GLOSS badge "5/5 hashes cracked. Elapsed: 0.003s" --type=success
sleep 0.3

echo ""
$GLOSS divider --style=heavy --gradient=matrix

# ── NETWORK SCAN ─────────────────────────────────────────────────
echo ""
$GLOSS badge "Launching nmap stealth scan (SYN -sS -O -A -T4)" --type=info
sleep 0.2
target_ip="$(randip)"
printf "${D}    Target: $target_ip${NC}\n\n"
$GLOSS table \
  "22/tcp=SSH open [$(randhex 4)]" \
  "80/tcp=HTTP open [$(randhex 4)]" \
  "443/tcp=HTTPS open [$(randhex 4)]" \
  "3306/tcp=MySQL open [$(randhex 4)]" \
  "6379/tcp=Redis open [NO AUTH]" \
  --style=single --gradient=matrix
echo ""
$GLOSS callout "Redis running with NO authentication. Jackpot." --type=error
sleep 0.4

echo ""
$GLOSS divider --style=heavy --gradient=matrix

# ── FAKE BITCOIN WALLET ───────────────────────────────────────────
echo ""
$GLOSS badge "Scanning blockchain for unprotected wallets..." --type=warning
sleep 0.3
printf "${D}\n"
for i in $(seq 1 6); do
  addr="1$(randhex 33 | tr 'a-f' 'A-F' | cut -c1-33)"
  bal="0.00000$(randint 1 9)"
  printf "  Wallet: ${Y}%s${NC}  Balance: ${G}%s BTC${NC}\n" "$addr" "$bal"
  sleep 0.2
done
printf "${NC}\n"
$GLOSS callout "Hot wallet found — Balance: 4.2069 BTC (~\$241,337)" --type=error
spinner "Draining wallet to cold storage..." 3

echo ""
$GLOSS divider --style=heavy --gradient=matrix

# ── FAKE AI / NEURAL NETWORK ──────────────────────────────────────
echo ""
$GLOSS "AI PWNING" --font=small --gradient=synthwave --border=none
$GLOSS badge "Initializing adversarial neural network..." --type=info
sleep 0.2
printf "${D}\n"
for epoch in $(seq 1 8); do
  loss=$(echo "scale=4; 2.5 / $epoch" | bc 2>/dev/null || echo "0.3125")
  acc=$(( 50 + epoch * 6 ))
  printf "  Epoch %d/8  loss=%-8s  acc=%d%%  val_acc=%d%%\n" \
    "$epoch" "$loss" "$acc" "$(( acc - 2 ))"
  sleep 0.3
done
printf "${NC}\n"
$GLOSS spark 50,56,62,68,74,80,86,92 --gradient=aurora
echo ""
$GLOSS badge "Model converged. Adversarial attack vector generated." --type=success
$GLOSS badge "Target AI confidence poisoned: 99.7% → 0.3%" --type=error
sleep 0.3

echo ""
$GLOSS divider --style=heavy --gradient=matrix

# ── SATELLITE UPLINK ─────────────────────────────────────────────
echo ""
$GLOSS badge "Hijacking satellite uplink — NORAD TLE acquired..." --type=warning
sleep 0.3
printf "${D}"
printf "  SATCOM-%d\n" "$(randint 1000 9999)"
printf "  TLE Line 1: 1 %dU 98067A   24%d.%08d  .%s  %s\n" \
  "$(randint 10000 99999)" "$(randint 100 365)" "$(randint 10000000 99999999)" "$(randhex 8)" "$(randhex 5)"
printf "  Ground track: LAT %d.%04d LON %d.%04d\n" \
  "$(randint -90 90)" "$(randint 1000 9999)" "$(randint -180 180)" "$(randint 1000 9999)"
printf "  Signal: -%d dBm  Freq: %d.%03d MHz\n" \
  "$(randint 40 120)" "$(randint 400 450)" "$(randint 100 999)"
printf "${NC}\n"
spinner "Uplink synchronized. Injecting payload..." 2

echo ""
$GLOSS divider --style=heavy --gradient=matrix

# ── EXPLOIT COMPILATION ───────────────────────────────────────────
echo ""
$GLOSS badge "Compiling zero-day exploit chain (CVE-2024-$(randint 10000 99999))" --type=warning
printf "${D}\n"
for f in heap_spray.c rop_gadget.c kernel_overwrite.c privesc.c rootkit_loader.c persistence.c; do
  printf "  gcc -O2 -m64 -fno-stack-protector -z execstack %s -o %s\n" "$f" "${f%.c}"
  sleep 0.25
done
printf "${NC}\n"
$GLOSS bar 100 --style=thin --label="Linking" --width=60 --gradient=fire
$GLOSS badge "Binary ready: ./$(randword)_exploit [57,344 bytes, stripped]" --type=success

echo ""
$GLOSS divider --style=heavy --gradient=matrix

# ── FILE EXFIL ──────────────────────────────────────────────────
echo ""
$GLOSS badge "Directory traversal + exfiltration:" --type=info
echo ""
$GLOSS list \
  ".ssh/id_rsa — EXFILTRATED (4096-bit RSA)" \
  ".ssh/authorized_keys — BACKDOOR INJECTED" \
  "/etc/shadow — CRACKED: 23/23 hashes" \
  ".aws/ — KEYS YOINKED (\$847,203 credit)" \
  ".kube/config — CLUSTER OWNED" \
  "/var/log/ — SHREDDED" \
  --style=arrow --gradient=fire
echo ""
spinner "Uploading to darknet dropzone via onion..." 3

echo ""
$GLOSS divider --style=heavy --gradient=matrix

# ── DDOS ─────────────────────────────────────────────────────────
echo ""
$GLOSS badge "Spinning up botnet — 12,847 nodes online" --type=error
sleep 0.3
printf "${D}\n"
for i in $(seq 1 12); do
  printf "  [node-%04d] %-18s → target: %-16s  %d req/s\n" \
    "$i" "$(randip)" "$(randip)" "$(randint 10000 99999)"
  sleep 0.08
done
printf "${NC}\n"
$GLOSS bar 100 --style=blocks --label="DDoS ramping" --width=60 --gradient=fire
$GLOSS badge "Target $(randip) DOWN — 0% uptime. CEO crying." --type=error
sleep 0.3

echo ""
$GLOSS divider --style=heavy --gradient=matrix

# ── HELLO NEO ────────────────────────────────────────────────────
echo ""
$GLOSS "HELLO NEO" --font=slant --gradient=rainbow --border=rounded
echo ""
printf "${LG}"
slow "    Wake up, Neo..."
sleep 0.8
slow "    The Matrix has you..."
sleep 0.8
slow "    Follow the white rabbit."
printf "${NC}\n"
sleep 1

echo ""
$GLOSS divider --style=heavy --gradient=matrix

# ── SELF-DESTRUCT ────────────────────────────────────────────────
echo ""
$GLOSS callout "Trace detected from $(randip) — FBI Cyber Division" --type=error
sleep 0.3
$GLOSS callout "Warrant issued — case #$(randint 100000 999999) — initiating self-destruct" --type=error
echo ""
for i in 10 9 8 7 6 5 4 3 2 1; do
  printf "${R}  SELF-DESTRUCT IN %2d SECONDS${NC}\n" "$i"
  sleep 0.5
done

echo ""
$GLOSS "BOOM" --font=banner --gradient=fire --border=none
sleep 0.5

echo ""
$GLOSS divider --style=double --gradient=aurora
echo ""
$GLOSS badge "Just kidding. Relax." --type=success
$GLOSS badge "No laws were broken. Probably." --type=success
$GLOSS badge "Have a great day, $(whoami)." --type=info
echo ""
