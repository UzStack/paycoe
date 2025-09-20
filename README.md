# Paycue docs

Assalomu Alaykum dasturdan foydalanishdan avval bularni o’qishingiz kerak loyiha qanday ishlaydi va qachon ishlatishingiz kerakligi haqida

Loyiha nomi paycue to’lovlarni avtomatlashtirish uchun open source dastur. Bu dastur yordamida siz to’lov tizimlariga integratsiya qilmasdan to’lovlarni avtomatlashtirishingiz mumkun.

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

> O’qishingiz shart: Telegram account ochilgan no’merda plastik karta sms xabarnoma yoqilgan bo’lishi shart
> 

# Quickstart

### O’rnatish

Githubdan oxirgi releaseni yuklab oling [download](https://github.com/JscorpTech/paycoe) `<arch>` o’rniga serveringizdagi arch yoziladi odatda `amd`

```bash
curl -o paycue -L https://github.com/UzStack/paycoe/releases/download/<version>/paycue-linux-<arch>
```

dastur uchun papka yaratishimiz kerak `/opt` papkasiga yaratishni maslahat beraman

```bash
mkdir -p /opt/paycue
```

va dasturni shu papkaga ko’chiring

```bash
mv ./paycue /opt/paycue
```

fayil uchun kerakli permissionlarni beramiz

```bash
sudo chmod +x ./paycue
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
PORT=10800
```

> Eslatma: .env fayildagi `CALLBACK_URL` juda muhum to’lov bajarilgandan keyin shu callback urlga malumotlarni yuboradi qaysi transaction bajarilganligi haqida
> 

### Botni sozlash

Keyingi navbat telegram botni sozlashimiz kerak [@HUMOcardbot](https://t.me/HUMOcardbot) ga kiring va botdagi ko’rsatmalarga amal qilib ro’yhatdan o’ting.

> To’lovlar uchun ishlatmoqchi bo’lgan kartangiz `💳 Kartalarni boshqarish` bo’limida mavjud kanligini tekshiring
> 

### Telegram accountni ulash

Telegram accountni dasturga ulash uchun bu commanddan foydalaning. Ko’rsatmalarga amal qiling

```bash
./paycue --telegram
```

### systemdni sozlash

Dastur doimiy ishlashi uchun systemd yordamida ishga tushuramiz 

yangi fayil yarating `/etc/systemd/system/paycue.service` 

```bash
[Unit]
Description="paycue service"
After=network.target

[Service]
User=root
Group=root
Type=simple
Restart=on-failure
RestartSec=5s
ExecStart=/opt/paycue/paycue
WorkingDirectory=/opt/paycue/

[Install]
WantedBy=multi-user.target
```

deyarli tayyor endi systemd ni  ishga tushursak bo’ldi

```bash
sudo systemctl enable --now paycue
```

dastur ishlayotganini tekshiring

```bash
sudo systemctl status paycue
```

# Muhum malumotlar

- To’lovdan avval transaction yaratasiz va dastur qaytargan miqdorda to’lov qilishini so’raysiz
- Transaction 30 daqiqa active qoladi keyin bekor qilinadi 30 daqiqadan keyingi to’langan to’lovlar tasdiqlanmaydi.
- Dastur ko’plab transactionlar bilan ishlay oladi lekin to’lov summasi farqi kattalshib ketishi mumkun masalan `10 ming` so’mlik `1000 ta` transactiondan keyin to’lov `11 ming` bo’lib ketadi buni oldini olish uchun bir nechta kartalardan foydalanishingiz mumkun dasturni bir nechta varintlarini turli accountlarga ulaysiz. (`buni hozirda qo’lda so’zlashingiz kerak  keyingi yangilanishlarda buni avtomatlashtiramiz`)
