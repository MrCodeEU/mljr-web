#!/usr/bin/env bash
# E2E smoke test for the newsletter app: boots the real binary against a
# scratch pb_data dir, drives 3 accounts through signup/login/groups/
# invites/answering, manipulates dates to walk the reminder -> grace ->
# close -> auto-disable lifecycle across real cron ticks, and verifies every
# outcome via the superuser REST API (email_log, notifications, etc.) rather
# than scraping HTML.
#
# Cron runs every 5 minutes in production (wall-clock aligned to */5 * * * *),
# but this script boots the server with NEWSLETTER_E2E_DEBUG=1, which adds a
# POST /debug/scan-now route that runs scheduler.RunScan synchronously and
# returns only once it's fully done — so the script drives the reminder ->
# grace -> close -> auto-disable lifecycle by calling scan_now() directly
# instead of waiting on real ticks.
#
# Usage: projects/newsletter/e2e/run.sh   (run from the repo root)
set -u
cd "$(dirname "$0")/../../.."

DATA_DIR=/tmp/nl-e2e
PORT=8097
BASE="http://127.0.0.1:$PORT"
SU_EMAIL="su-e2e@example.com"
SU_PASS="su-e2e-pass-1234"
LOG="$DATA_DIR.log"

PASS=0
FAIL=0
SERVER_PID=""

# ---- output / assertions ---------------------------------------------------

note() { printf '\n\033[1;36m== %s\033[0m\n' "$*"; }

assert_eq() {
  local desc="$1" expected="$2" actual="$3"
  if [ "$expected" = "$actual" ]; then
    PASS=$((PASS + 1)); printf '  \033[32m✓\033[0m %s\n' "$desc"
  else
    FAIL=$((FAIL + 1)); printf '  \033[31m✗\033[0m %s (expected %q, got %q)\n' "$desc" "$expected" "$actual"
  fi
}

assert_true() {
  local desc="$1" cond="$2"
  case "$cond" in
    true) PASS=$((PASS + 1)); printf '  \033[32m✓\033[0m %s\n' "$desc" ;;
    [1-9]*) PASS=$((PASS + 1)); printf '  \033[32m✓\033[0m %s\n' "$desc" ;;
    *) FAIL=$((FAIL + 1)); printf '  \033[31m✗\033[0m %s (got %q)\n' "$desc" "$cond" ;;
  esac
}

cleanup() {
  local code=$?
  if [ -n "$SERVER_PID" ] && kill -0 "$SERVER_PID" 2>/dev/null; then
    kill "$SERVER_PID" 2>/dev/null
    wait "$SERVER_PID" 2>/dev/null
  fi
  if [ "$FAIL" -gt 0 ] || [ "$code" -ne 0 ]; then
    echo "--- server log tail (left $DATA_DIR for inspection) ---"
    tail -n 60 "$LOG" 2>/dev/null
  else
    rm -rf "$DATA_DIR" "$LOG"
  fi
  echo
  echo "PASS=$PASS FAIL=$FAIL"
  [ "$FAIL" -eq 0 ]
}
trap cleanup EXIT

# ---- superuser REST helpers --------------------------------------------------

SU_TOKEN=""
pb_auth() {
  SU_TOKEN=$(curl -s -X POST "$BASE/api/collections/_superusers/auth-with-password" \
    -H 'Content-Type: application/json' \
    -d "{\"identity\":\"$SU_EMAIL\",\"password\":\"$SU_PASS\"}" | jq -r '.token')
  [ -n "$SU_TOKEN" ] && [ "$SU_TOKEN" != "null" ]
}

# pb_list collection filter -> JSON array of matching records (admin-only read)
pb_list() {
  local col="$1" filter="$2"
  curl -s -G "$BASE/api/collections/$col/records" \
    -H "Authorization: Bearer $SU_TOKEN" \
    --data-urlencode "filter=$filter" --data-urlencode "perPage=200" \
    | jq -c '.items'
}

# pb_get_one collection filter -> first matching record JSON, or "null"
pb_get_one() { pb_list "$1" "$2" | jq -c '.[0] // null'; }

# pb_create collection json_body -> created record JSON
pb_create() {
  curl -s -X POST "$BASE/api/collections/$1/records" \
    -H "Authorization: Bearer $SU_TOKEN" -H 'Content-Type: application/json' \
    -d "$2"
}

# pb_patch collection id json_body -> updated record JSON
pb_patch() {
  curl -s -X PATCH "$BASE/api/collections/$1/records/$2" \
    -H "Authorization: Bearer $SU_TOKEN" -H 'Content-Type: application/json' \
    -d "$3"
}

# iso_offset seconds -> PocketBase datetime string for now+seconds (negative = past)
iso_offset() { date -u -d "@$(($(date +%s) + $1))" +"%Y-%m-%d %H:%M:%S.000Z"; }

# ---- app HTTP helpers --------------------------------------------------------

# post_form jar path field=val... -> prints http status code
post_form() {
  local jar="$1" path="$2"; shift 2
  local args=(-s -o /dev/null -w '%{http_code}' -b "$jar" -c "$jar" -X POST "$BASE$path")
  for kv in "$@"; do args+=(--data-urlencode "$kv"); done
  curl "${args[@]}"
}

signup() { post_form "$1" "/signup" "email=$2" "password=$3" "name=$4" ${5:+"invite=$5"}; }
login()  { post_form "$1" "/login"  "email=$2" "password=$3" ${4:+"invite=$4"}; }

# scan_now triggers an immediate scheduler pass via the debug-only endpoint
# (NEWSLETTER_E2E_DEBUG=1) and blocks until it's fully done, replacing the
# old "wait up to 360s for a real cron tick" pattern.
scan_now() { curl -s -o /dev/null -X POST "$BASE/debug/scan-now"; }

# ---- setup --------------------------------------------------------------------

note "build + bootstrap scratch instance"
rm -rf "$DATA_DIR"
make build PROJECT=newsletter >/dev/null || { echo "build failed"; exit 1; }
./bin/newsletter superuser upsert "$SU_EMAIL" "$SU_PASS" --dir="$DATA_DIR" >/dev/null || { echo "superuser bootstrap failed"; exit 1; }

NEWSLETTER_E2E_DEBUG=1 ./bin/newsletter serve --http=127.0.0.1:$PORT --dir="$DATA_DIR" >"$LOG" 2>&1 &
SERVER_PID=$!
for i in $(seq 1 60); do
  curl -s -o /dev/null "$BASE/healthz" && break
  sleep 1
done
curl -s "$BASE/healthz" | grep -q ok || { echo "server never came up"; exit 1; }
pb_auth || { echo "superuser auth failed"; exit 1; }
echo "server up, pid=$SERVER_PID, data dir=$DATA_DIR"

TINY_PNG=$(mktemp --suffix=.png)
base64 -d >"$TINY_PNG" <<'EOF'
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII=
EOF

# ---- step 1+2: signup + login --------------------------------------------------

note "step 1-2: signup 3 accounts, then re-login each to prove login independently works"
JAR_A=$(mktemp); JAR_B=$(mktemp); JAR_C=$(mktemp)
EMAIL_A="alice@example.com"; EMAIL_B="bob@example.com"; EMAIL_C="${E2E_REAL_RECIPIENT:-carol@example.com}"
PASS_AB="password123"

assert_eq "signup A" 303 "$(signup "$JAR_A" "$EMAIL_A" "$PASS_AB" "Alice")"
assert_eq "signup B" 303 "$(signup "$JAR_B" "$EMAIL_B" "$PASS_AB" "Bob")"
assert_eq "signup C" 303 "$(signup "$JAR_C" "$EMAIL_C" "$PASS_AB" "Carol")"

: >"$JAR_A"; : >"$JAR_B"; : >"$JAR_C" # drop the signup session, force a real login
assert_eq "login A" 303 "$(login "$JAR_A" "$EMAIL_A" "$PASS_AB")"
assert_eq "login B" 303 "$(login "$JAR_B" "$EMAIL_B" "$PASS_AB")"
assert_eq "login C" 303 "$(login "$JAR_C" "$EMAIL_C" "$PASS_AB")"

USER_A=$(pb_get_one users "email='$EMAIL_A'" | jq -r '.id')
USER_B=$(pb_get_one users "email='$EMAIL_B'" | jq -r '.id')
USER_C=$(pb_get_one users "email='$EMAIL_C'" | jq -r '.id')
assert_true "all 3 users exist" "$([ -n "$USER_A" ] && [ -n "$USER_B" ] && [ -n "$USER_C" ] && echo true)"

# ---- step 3+4: groups, invites, joins ------------------------------------------

note "step 3-4: A creates groups, invites B+C, they join"

create_group_as_a() {
  post_form "$JAR_A" "/groups" "name=$1" >/dev/null
  pb_get_one groups "name='$1'" | jq -r '.id + " " + .slug'
}
invite_and_accept() {
  local slug="$1" jar="$2" email="$3"
  post_form "$JAR_A" "/g/$slug/invites" "email=$email" "role=member" >/dev/null
  local token
  token=$(pb_get_one group_invites "group.slug='$slug' && email='$email'" | jq -r '.token')
  assert_true "invite token issued for $email/$slug" "$([ -n "$token" ] && [ "$token" != "null" ] && echo true)"
  assert_eq "accept invite $email/$slug" 303 "$(post_form "$jar" "/invites/$token/accept")"
}

read -r G1_ID G1_SLUG <<<"$(create_group_as_a "Movie Night")"
read -r G2_ID G2_SLUG <<<"$(create_group_as_a "Book Club")"
read -r GHOST_ID GHOST_SLUG <<<"$(create_group_as_a "Ghost Crew")"
read -r CURATE_ID CURATE_SLUG <<<"$(create_group_as_a "Curate Crew")"
assert_true "4 groups created" "$([ -n "$G1_ID" ] && [ -n "$G2_ID" ] && [ -n "$GHOST_ID" ] && [ -n "$CURATE_ID" ] && echo true)"

invite_and_accept "$G1_SLUG" "$JAR_B" "$EMAIL_B"
invite_and_accept "$G1_SLUG" "$JAR_C" "$EMAIL_C"
invite_and_accept "$G2_SLUG" "$JAR_B" "$EMAIL_B"
invite_and_accept "$G2_SLUG" "$JAR_C" "$EMAIL_C"
invite_and_accept "$CURATE_SLUG" "$JAR_B" "$EMAIL_B"

assert_eq "G1 has 3 members" 3 "$(pb_list group_memberships "group='$G1_ID'" | jq 'length')"
assert_eq "G2 has 3 members" 3 "$(pb_list group_memberships "group='$G2_ID'" | jq 'length')"

# ---- step 5: answer flow + grace-period straggler -------------------------------

note "step 5: A opens a G1 edition, A+B answer everything, C deliberately doesn't (C is the real-recipient slot, so the reminder/grace emails land in a real inbox)"
assert_eq "A opens G1 edition" 303 "$(post_form "$JAR_A" "/g/$G1_SLUG/editions")"
G1_EDITION=$(pb_get_one newsletter_editions "group='$G1_ID'" | jq -r '.id')
assert_true "G1 edition created+open" "$([ -n "$G1_EDITION" ] && echo true)"

EQS=$(pb_list edition_questions "edition='$G1_EDITION'")
QCOUNT=$(echo "$EQS" | jq 'length')
assert_eq "G1 edition has all 30 seeded global questions" 30 "$QCOUNT"

# answer_args jar question_id type options_json -> populated curl -F args via stdout, one per line
answer_args() {
  local qid="$1" type="$2" opts="$3"
  case "$type" in
    text) echo "-F"; echo "q_$qid=E2E test answer" ;;
    single_select|emoji_reaction)
      echo "-F"; echo "q_$qid=$(echo "$opts" | jq -r '.[0]')" ;;
    multi_select)
      for v in $(echo "$opts" | jq -r '.[0:2][]'); do echo "-F"; echo "q_$qid[]=$v"; done ;;
    rating) echo "-F"; echo "q_$qid=4" ;;
    image) echo "-F"; echo "f_$qid=@$TINY_PNG;type=image/png" ;;
    yes_no) echo "-F"; echo "q_$qid=true" ;;
    scale) echo "-F"; echo "q_$qid=7" ;;
    number) echo "-F"; echo "q_$qid=3" ;;
    date) echo "-F"; echo "q_$qid=2026-01-15" ;;
    color_pick) echo "-F"; echo "q_$qid=#8b5cf6" ;;
  esac
}

# answer_all jar who edition_id eqs_json -> submits an answer for every
# question in eqs_json against edition_id, as jar's user.
answer_all() {
  local jar="$1" who="$2" edition_id="$3" eqs="$4"
  local curl_args=(-s -o /dev/null -w '%{http_code}' -b "$jar" -c "$jar"
    -X POST "$BASE/g/$G1_SLUG/editions/$edition_id")
  while IFS= read -r qid; do
    local type opts
    type=$(pb_list question_bank "id='$qid'" | jq -r '.[0].type')
    opts=$(pb_list question_bank "id='$qid'" | jq -c '.[0].options // []')
    while IFS= read -r arg; do curl_args+=("$arg"); done < <(answer_args "$qid" "$type" "$opts")
  done < <(echo "$eqs" | jq -r '.[].question')
  local code
  code=$(curl "${curl_args[@]}")
  assert_eq "$who submits all answers for $edition_id" 303 "$code"
}

answer_all "$JAR_A" "A" "$G1_EDITION" "$EQS"
answer_all "$JAR_B" "B" "$G1_EDITION" "$EQS"

ANSWERED_A=$(pb_list answers "edition='$G1_EDITION' && user='$USER_A' && skipped=false" | jq 'length')
ANSWERED_B=$(pb_list answers "edition='$G1_EDITION' && user='$USER_B' && skipped=false" | jq 'length')
assert_eq "A answered all 30" 30 "$ANSWERED_A"
assert_eq "B answered all 30" 30 "$ANSWERED_B"
assert_eq "C answered nothing yet" 0 "$(pb_list answers "edition='$G1_EDITION' && user='$USER_C' && skipped=false" | jq 'length')"

# ---- step 7: settings/profile edge cases ----------------------------------------

note "step 7: profile rename, group settings rename, C leaves G2"
assert_eq "A renames profile" 303 "$(post_form "$JAR_A" "/profile" "name=Alice Renamed")"
assert_eq "profile name persisted" "Alice Renamed" "$(curl -s "$BASE/api/collections/users/records/$USER_A" -H "Authorization: Bearer $SU_TOKEN" | jq -r '.name')"

assert_eq "A renames + retunes G1 settings" 303 "$(post_form "$JAR_A" "/g/$G1_SLUG/settings" \
  "name=Movie Night Weekly" "schedule_period=weekly" "reminder_lead_hours=1" "grace_period_hours=1")"
G1_AFTER=$(curl -s "$BASE/api/collections/groups/records/$G1_ID" -H "Authorization: Bearer $SU_TOKEN")
assert_eq "G1 name updated" "Movie Night Weekly" "$(echo "$G1_AFTER" | jq -r '.name')"
assert_eq "G1 grace_period_hours updated" 1 "$(echo "$G1_AFTER" | jq -r '.grace_period_hours')"

assert_eq "C leaves G2" 303 "$(post_form "$JAR_C" "/g/$G2_SLUG/leave")"
assert_eq "C's G2 membership gone" "null" "$(pb_get_one group_memberships "group='$G2_ID' && user='$USER_C'")"

# ---- step 7b: full profile attribute coverage -----------------------------------

note "step 7b: profile birthday/favorite_animal/favorite_food/favorite_color + avatar upload"
assert_eq "B sets birthday/animal/food/color" 303 "$(post_form "$JAR_B" "/profile" \
  "name=Bob" "birthday=1990-07-15" "favorite_animal=Otter" "favorite_food=Pizza" "favorite_color=violet")"
USER_B_AFTER=$(curl -s "$BASE/api/collections/users/records/$USER_B" -H "Authorization: Bearer $SU_TOKEN")
assert_eq "B birthday persisted" "1990-07-15" "$(echo "$USER_B_AFTER" | jq -r '.birthday' | cut -c1-10)"
assert_eq "B favorite_animal persisted" "Otter" "$(echo "$USER_B_AFTER" | jq -r '.favorite_animal')"
assert_eq "B favorite_food persisted" "Pizza" "$(echo "$USER_B_AFTER" | jq -r '.favorite_food')"
assert_eq "B favorite_color persisted" "violet" "$(echo "$USER_B_AFTER" | jq -r '.favorite_color')"

assert_eq "B submits invalid favorite_color, request rejected" 303 "$(post_form "$JAR_B" "/profile" \
  "name=Bob" "favorite_color=notacolor")"
assert_eq "invalid favorite_color did not overwrite stored value" "violet" \
  "$(curl -s "$BASE/api/collections/users/records/$USER_B" -H "Authorization: Bearer $SU_TOKEN" | jq -r '.favorite_color')"

AVATAR_CODE=$(curl -s -o /dev/null -w '%{http_code}' -b "$JAR_B" -c "$JAR_B" \
  -X POST "$BASE/profile/avatar" -F "avatar=@$TINY_PNG;type=image/png")
assert_eq "B uploads avatar" 303 "$AVATAR_CODE"
AVATAR_FIELD=$(curl -s "$BASE/api/collections/users/records/$USER_B" -H "Authorization: Bearer $SU_TOKEN" | jq -r '.avatar')
assert_true "B's avatar field populated" "$([ -n "$AVATAR_FIELD" ] && [ "$AVATAR_FIELD" != "null" ] && echo true)"

# ---- step 7c: group settings cover every schedule_period option ----------------

note "step 7c: G2 settings cycle through every schedule_period option"
for period in weekly biweekly monthly quarterly; do
  CODE=$(post_form "$JAR_A" "/g/$G2_SLUG/settings" \
    "name=Book Club" "schedule_period=$period" "reminder_lead_hours=24" "grace_period_hours=12")
  assert_eq "G2 settings save with schedule_period=$period" 303 "$CODE"
  assert_eq "G2 schedule_period=$period persisted" "$period" \
    "$(curl -s "$BASE/api/collections/groups/records/$G2_ID" -H "Authorization: Bearer $SU_TOKEN" | jq -r '.schedule_period')"
done

# ---- ghost group: pre-seed 3 dead editions before any tick fires ----------------

note "seeding Ghost Crew with 3 already-expired zero-answer editions (auto-disable test)"
for i in 1 2 3; do
  pb_create newsletter_editions "{\"group\":\"$GHOST_ID\",\"status\":\"open\",\"opens_at\":\"$(iso_offset -3600)\",\"closes_at\":\"$(iso_offset -1800)\",\"grace_until\":\"$(iso_offset -60)\"}" >/dev/null
done
assert_eq "ghost group has 3 pre-seeded editions" 3 "$(pb_list newsletter_editions "group='$GHOST_ID'" | jq 'length')"

# ---- curate group: pre-seed a scheduled edition + candidate questions ----------

note "seeding Curate Crew's next edition + adding a custom question, then curating it down"
assert_eq "A adds a custom question to Curate Crew" 303 "$(post_form "$JAR_A" "/g/$CURATE_SLUG/questions" \
  "prompt=Custom curate question" "type=text")"
CUSTOM_Q=$(pb_get_one question_bank "group='$CURATE_ID' && prompt='Custom curate question'" | jq -r '.id')
assert_true "custom question created" "$([ -n "$CUSTOM_Q" ] && [ "$CUSTOM_Q" != "null" ] && echo true)"

CANDIDATES=$(pb_list question_bank "is_active=true && (scope='global' || (scope='group' && group='$CURATE_ID'))")
CAND_COUNT=$(echo "$CANDIDATES" | jq 'length')
assert_eq "curate candidate pool is 30 global + 1 custom" 31 "$CAND_COUNT"

CURATE_EDITION=$(pb_create newsletter_editions "{\"group\":\"$CURATE_ID\",\"status\":\"scheduled\",\"opens_at\":\"$(iso_offset -60)\"}" | jq -r '.id')
i=0
DROPPED_Q=""
while IFS= read -r qid; do
  if [ -z "$DROPPED_Q" ]; then DROPPED_Q="$qid"; continue; fi # drop the first candidate
  pb_create edition_questions "{\"edition\":\"$CURATE_EDITION\",\"question\":\"$qid\",\"order\":$i}" >/dev/null
  i=$((i + 1))
done < <(echo "$CANDIDATES" | jq -r '.[].id')
assert_eq "curate edition seeded with 30 of 31 candidates" 30 "$(pb_list edition_questions "edition='$CURATE_EDITION'" | jq 'length')"

# curate via the real admin UI route: re-add the dropped question, drop a different one instead
KEEP_ARGS=(-s -o /dev/null -w '%{http_code}' -b "$JAR_A" -c "$JAR_A" -X POST "$BASE/g/$CURATE_SLUG/editions/$CURATE_EDITION/questions")
SURVIVOR_DROPPED=""
order=0
while IFS= read -r qid; do
  if [ -z "$SURVIVOR_DROPPED" ] && [ "$qid" != "$DROPPED_Q" ]; then SURVIVOR_DROPPED="$qid"; continue; fi
  KEEP_ARGS+=(--data-urlencode "q_$qid=1" --data-urlencode "order_$qid=$order")
  order=$((order + 1))
done < <(echo "$CANDIDATES" | jq -r '.[].id')
CODE=$(curl "${KEEP_ARGS[@]}")
assert_eq "A curates the edition's question set" 303 "$CODE"

CURATED_SET=$(pb_list edition_questions "edition='$CURATE_EDITION'")
assert_eq "curated edition keeps 30 questions" 30 "$(echo "$CURATED_SET" | jq 'length')"
assert_true "previously-dropped question is back in" "$(echo "$CURATED_SET" | jq --arg q "$DROPPED_Q" '[.[].question] | index($q) != null')"
assert_true "newly-dropped question is now out" "$(echo "$CURATED_SET" | jq --arg q "$SURVIVOR_DROPPED" '[.[].question] | index($q) == null')"

# ---- tick A: reminders + curate-edition opens + ghost group disables -----------

note "backdating G1's reminder_at, then forcing tick A"
pb_patch newsletter_editions "$G1_EDITION" "{\"reminder_at\":\"$(iso_offset -60)\"}" >/dev/null
scan_now

assert_true "G1 reminder email logged to C" \
  "$([ "$(pb_get_one email_log "kind='reminder' && dedupe_key='reminder:$G1_EDITION:$USER_C'")" != "null" ] && echo true)"
assert_eq "G1 reminder logged for C only" 1 "$(pb_list email_log "kind='reminder' && dedupe_key~'reminder:$G1_EDITION:'" | jq 'length')"

assert_eq "curate edition opened" "open" "$(curl -s "$BASE/api/collections/newsletter_editions/records/$CURATE_EDITION" -H "Authorization: Bearer $SU_TOKEN" | jq -r '.status')"
BANSWER_HTML=$(curl -s -b "$JAR_B" "$BASE/g/$CURATE_SLUG/editions/$CURATE_EDITION")
assert_true "B's answer page shows the surviving custom question" "$(echo "$BANSWER_HTML" | grep -qF 'Custom curate question' && echo true || echo false)"

# scan_now blocks until RunScan (sendReminders -> ... -> closeEditions) has
# fully returned, so the ghost group's disable is already committed here —
# no separate poll/wait needed like with the old real-cron-tick approach.
GHOST_AFTER=$(curl -s "$BASE/api/collections/groups/records/$GHOST_ID" -H "Authorization: Bearer $SU_TOKEN")
assert_eq "ghost group streak hit 3" 3 "$(echo "$GHOST_AFTER" | jq -r '.consecutive_unanswered_editions')"
assert_eq "ghost group auto-disabled" "disabled" "$(echo "$GHOST_AFTER" | jq -r '.status')"
assert_eq "ghost group's 3 editions all closed" 3 "$(pb_list newsletter_editions "group='$GHOST_ID' && status='sent'" | jq 'length')"
assert_true "ghost group got a group_disabled notification" "$([ "$(pb_list notifications "group='$GHOST_ID' && kind='group_disabled'" | jq 'length')" -gt 0 ] && echo true)"

# ---- tick B: grace reminders + ghost group stays inert --------------------------

note "backdating G1's closes_at, then forcing tick B"
pb_patch newsletter_editions "$G1_EDITION" "{\"closes_at\":\"$(iso_offset -60)\"}" >/dev/null
scan_now

assert_true "G1 grace reminder logged to C" \
  "$([ "$(pb_get_one email_log "kind='grace_reminder' && dedupe_key='grace_reminder:$G1_EDITION:$USER_C'")" != "null" ] && echo true)"
assert_eq "grace reminder went to C only, not A/B" 1 "$(pb_list email_log "kind='grace_reminder' && dedupe_key~'grace_reminder:$G1_EDITION:'" | jq 'length')"
assert_eq "ghost group still has exactly 3 editions (no resurrection)" 3 "$(pb_list newsletter_editions "group='$GHOST_ID'" | jq 'length')"

# ---- tick C: G1 closes, A's unanswered questions get skipped -------------------

note "backdating G1's grace_until, then forcing tick C"
pb_patch newsletter_editions "$G1_EDITION" "{\"grace_until\":\"$(iso_offset -60)\"}" >/dev/null
scan_now

assert_eq "G1 edition sent" "sent" "$(curl -s "$BASE/api/collections/newsletter_editions/records/$G1_EDITION" -H "Authorization: Bearer $SU_TOKEN" | jq -r .status)"
assert_eq "G1 edition_sent email logged, one per member" 3 "$(pb_list email_log "kind='edition_sent' && dedupe_key~'send:$G1_EDITION:'" | jq 'length')"
assert_eq "C's 30 questions are now marked skipped" 30 "$(pb_list answers "edition='$G1_EDITION' && user='$USER_C' && skipped=true" | jq 'length')"
assert_eq "ghost group still disabled, still 3 editions (createDueEditions skips it)" 3 "$(pb_list newsletter_editions "group='$GHOST_ID'" | jq 'length')"

# ---- language switch: C's whole reminder->grace->sent cycle re-runs in German --
#
# logMailer (scheduler/mailer.go) logs "to=... subject=..." to $LOG since no
# SMTP_HOST is configured for e2e. reminder_subject/grace_subject differ by
# distinctive English/German words; sent_subject differs only by German noun
# capitalization ("Newsletter" vs "newsletter") - still a real, language-
# accurate signal, not a contrived one.

note "C's 1st-cycle mails were English; now switching C to German and re-running the cycle on a 2nd G1 edition"

assert_eq "C's 1st reminder was English" 1 "$(grep -F "$EMAIL_C" "$LOG" | grep -c "forget")"
assert_eq "C's 1st grace reminder was English" 1 "$(grep -F "$EMAIL_C" "$LOG" | grep -c "last call")"
assert_eq "C's 1st edition_sent was English" 1 "$(grep -F "$EMAIL_C" "$LOG" | grep -c " newsletter —")"

assert_eq "C sets language=de" 303 "$(post_form "$JAR_C" "/profile" "name=Carol" "language=de")"
assert_eq "C's language persisted as de" "de" \
  "$(curl -s "$BASE/api/collections/users/records/$USER_C" -H "Authorization: Bearer $SU_TOKEN" | jq -r '.language')"

assert_eq "A opens 2nd G1 edition" 303 "$(post_form "$JAR_A" "/g/$G1_SLUG/editions")"
G1_EDITION2=$(pb_list newsletter_editions "group='$G1_ID' && status='open'" | jq -r '.[0].id')
assert_true "2nd G1 edition created+open" "$([ -n "$G1_EDITION2" ] && [ "$G1_EDITION2" != "null" ] && echo true)"

EQS2=$(pb_list edition_questions "edition='$G1_EDITION2'")
answer_all "$JAR_A" "A" "$G1_EDITION2" "$EQS2"
answer_all "$JAR_B" "B" "$G1_EDITION2" "$EQS2"
assert_eq "C answered nothing on 2nd edition yet" 0 "$(pb_list answers "edition='$G1_EDITION2' && user='$USER_C' && skipped=false" | jq 'length')"

note "backdating 2nd G1 edition's reminder_at, then forcing tick A"
pb_patch newsletter_editions "$G1_EDITION2" "{\"reminder_at\":\"$(iso_offset -60)\"}" >/dev/null
scan_now
assert_true "2nd G1 reminder logged to C" \
  "$([ "$(pb_get_one email_log "kind='reminder' && dedupe_key='reminder:$G1_EDITION2:$USER_C'")" != "null" ] && echo true)"
assert_eq "2nd G1 reminder logged for C only" 1 "$(pb_list email_log "kind='reminder' && dedupe_key~'reminder:$G1_EDITION2:'" | jq 'length')"
assert_eq "C's 2nd reminder was German" 1 "$(grep -F "$EMAIL_C" "$LOG" | grep -c "Vergiss nicht")"

note "backdating 2nd G1 edition's closes_at, then forcing tick B"
pb_patch newsletter_editions "$G1_EDITION2" "{\"closes_at\":\"$(iso_offset -60)\"}" >/dev/null
scan_now
assert_true "2nd G1 grace reminder logged to C" \
  "$([ "$(pb_get_one email_log "kind='grace_reminder' && dedupe_key='grace_reminder:$G1_EDITION2:$USER_C'")" != "null" ] && echo true)"
assert_eq "2nd G1 grace reminder went to C only" 1 "$(pb_list email_log "kind='grace_reminder' && dedupe_key~'grace_reminder:$G1_EDITION2:'" | jq 'length')"
assert_eq "C's 2nd grace reminder was German" 1 "$(grep -F "$EMAIL_C" "$LOG" | grep -c "letzte Chance")"

note "backdating 2nd G1 edition's grace_until, then forcing tick C"
pb_patch newsletter_editions "$G1_EDITION2" "{\"grace_until\":\"$(iso_offset -60)\"}" >/dev/null
scan_now
assert_eq "2nd G1 edition sent" "sent" "$(curl -s "$BASE/api/collections/newsletter_editions/records/$G1_EDITION2" -H "Authorization: Bearer $SU_TOKEN" | jq -r .status)"
assert_eq "2nd G1 edition_sent email logged, one per member" 3 "$(pb_list email_log "kind='edition_sent' && dedupe_key~'send:$G1_EDITION2:'" | jq 'length')"
assert_eq "C's 2nd edition_sent was German" 1 "$(grep -F "$EMAIL_C" "$LOG" | grep -c " Newsletter —")"

# ---- idempotency: re-run the same window once more, expect no duplicates ------

note "idempotency check: confirm nothing double-sent if a tick re-processes the same state"
EMAIL_COUNT_BEFORE=$(pb_list email_log "dedupe_key~'$G1_EDITION'" | jq 'length')
scan_now # re-process the same already-closed state once more
EMAIL_COUNT_AFTER=$(pb_list email_log "dedupe_key~'$G1_EDITION'" | jq 'length')
assert_eq "no duplicate G1 email_log rows after re-processing" "$EMAIL_COUNT_BEFORE" "$EMAIL_COUNT_AFTER"

echo
echo "all steps executed."
