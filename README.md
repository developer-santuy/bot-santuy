# Bot developer santuy

## Create your bot
- Go to `@BotFather` in telegram to create your bot
- Rename `.env-example` with `.env`
- Put your bot token into `.env` file at `BOT_TOKEN`

## Local development
Use `ngrok`, you can download it [here](https://ngrok.com/download)  

### Run
Run ngrok on port 8000 or any port you want
```bash
ngrok http 8000
```

### Set webhook
You must activate the webhook first.

Run
```bash
curl -F "url=<YOUR_NGROK_URL>"  https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setWebhook
```

Or simply acccess
[https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setWebhook?url=<YOUR_NGROK_URL>](https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setWebhook?url=<YOUR_NGROK_URL>)


## Deploy in prod
The step is barely same, you dont need to use ngrok anymore since you using your own server

### Note
Your server must run over `HTTPS`