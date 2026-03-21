#!/bin/bash

TELEGRAM_BOT_API_TOKEN="8705197059:AAHU625kb6LbJaNy2J0Sc4IoSMJbX16lLtg"
BASE_URL=https://nai9yhaep0.execute-api.us-east-1.amazonaws.com/webhook
SECRET_TOKEN=abcdefghijklmnopqrstvwxyz1234567890
PAYLOAD="{
  \"url\": \"$BASE_URL\",
  \"secret_token\": \"$SECRET_TOKEN\",
  \"drop_pending_updates\": true
}
"

curl -H 'content-type: application/json' -XPOST --data "$PAYLOAD" -sS https://api.telegram.org/bot$TELEGRAM_BOT_API_TOKEN/setWebhook
