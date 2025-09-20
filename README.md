# Paycoe docs

Assalomu Alaykum dasturdan foydalanishdan avval bularni o’qishingiz kerak loyiha qanday ishlaydi va qachon ishlatishingiz kerakligi haqida

Loyiha nomi paycoe to’lovlarni avtomatlashtirish uchun open source dastur. Bu dastur yordamida siz to’lov tizimlariga integratsiya qilmasdan to’lovlarni avtomatlashtirishingiz mumkun.

`Dasturni kimlar ishlatishi kerak?` agar siz yuridik shaxs bo’lmasangiz lekin to’lovlarni avtomatlashtirmoqchi bo’lsangiz va foydalanuvchilaringiz ko’p bo’lmasa 10 ming dan kam bo’lsa bu dastur aynan siz uchun 

`Loyiha qanday ishlaydi?` yangi to’lov yaratish uchun yangi transaction yaratasiz masalan sizga `10 ming` so’m to’lov kerak 10 ming `amount` yuborasiz keyin dastur sizga avtomatik hozir active bo’lmagan summada amount qaytaradi masalan `1025 so’m` foydalanuvchidan shuncha pul to’lashini so’raysiz. va agrda sizning kartangizga berilgan summada pul tushsa dastur sizga xabar beradi api `webhook` yordamida

> Dasturdan foydalanishda savollaringiz bo’lsa telegram orqali [@Azamov_Samandar](https://t.me/Azamov_Samandar) ga yozishingiz mumkun
> 

### Kerakli narsalar

1. Telegram account
2. Humo telegram bot
3. Humo plastik kartasi
4. Server
5. Redis

`Nega Telegram account va humo kerak?` chunki dastur Humoning rasmiy botidan malumot olib ishlaydi. Humo kartaga pul tushganda humo telegram bot orqali xabar yuboradi dastur esa buni olib qayta ishlaydi.

# Quickstart

### Yuklab olish

Githubdan oxirgi releaseni yuklab oling [download](https://github.com/JscorpTech/paycoe)

```bash
curl -O https://github.com/JscorpTech/paycoe/archive/refs/tags/v1.0.0.zip
```

dastur uchun papka yaratishimiz kerak `/opt` papkasiga yaratishni maslahat beraman

```bash
mkdir -p /opt/paycoe
```

va dasturni shu papkaga ko’chiring

```bash
mv v1.0.0.zip /opt/paycoe
```

endi shu papkada `.env` fayil yaratishimiz kerak api hash va api keyni [my.telegram.org](http://my.telegram.org) saytidan olishingiz mumkun 

```bash
APP_ID=<app_id>
APP_HASH=<api_hash>
TG_PHONE=+998943990509
SESSION_DIR="sessions"
CALLBACK_URL=https://example.com
REDIS_ADDR=127.0.0.1:6379
WATCH_BOT_USERNAME=JscorpTechAdmin
WORKERS=10
```

test uchun dasturni ishga tushurib ko’ring

```bash
./paycoe
```

### systemdni sozlash

Dastur doimiy ishlashi uchun systemd yordamida ishga tushuramiz 

yangi fayil yarating `/etc/systemd/system/paycoe.service` 

```bash
[Unit]
Description="Paycoe service"
After=network.target

[Service]
User=root
Group=root
Type=simple
Restart=on-failure
RestartSec=5s
ExecStart=/opt/paycoe/paycoe
WorkingDirectory=/opt/paycoe/paycoe

[Install]
WantedBy=multi-user.target
```

deyarli tayyor endi systemd ni  ishga tushursak bo’ldi

```bash
sudo systemctl enable --now paycoe
```

dastur ishlayotganini tekshiring

```bash
sudo systemctl status paycoe
```
