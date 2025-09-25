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

Githubdan oxirgi releaseni yuklab oling [download](https://github.com/UzStack/paycoe) `<arch>` o’rniga serveringizdagi arch yoziladi odatda `amd`

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
APP_HASH=<app_hash>
TG_PHONE=<you_phone_number>
SESSION_DIR="sessions"
REDIS_ADDR=127.0.0.1:6379
WORKERS=10
WEBHOOK_URL=http://127.0.0.1:10800/health/
WATCH_ID=856254490
PORT=10800
DEBUG=true
LIMIT=100
```

> Eslatma: .env fayildagi `WEBHOOK_URL` juda muhum to’lov bajarilgandan keyin shu callback urlga malumotlarni yuboradi qaysi transaction bajarilganligi haqida
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

# Integratsiya

### Transaction yaratish

Request example

```bash
curl --request POST \
  --url http://<host>:10800/create/transaction/ \
  --header 'authorization: Basic OHFhS2ZhS3AyU19JZGZyNUlaU3dqeTFtSFFnYTpWMFl5VFZUMXkyQkRJUWVnVFdLYTI3bUtFU29h' \
  --header 'content-type: application/json' \
  --data '{
  "amount": 20000
}'
```

Post data

| amount | To’lov miqdori |
| --- | --- |

Success response misol

```json
{
  "status": true,
  "data": {
    "amount": 20000,
    "transaction_id": "622ea789-5b4c-4e6a-a76b-415ac144eb34"
  }
}
```

Error response misol

```json
{
  "status": false,
  "data": {
    "detail": "Amount must be less than 100"
  }
}
```

### Webhook

To’lov bajarilganda yoki bekor qilinganda dastur siz kiritgan callback urlga malumotlarni yuboradi. Ikkita asosiy  action mavjud cancel va confirm

```json
# to'lov bajarilganda
{
	"action": "confirm",
	"amount": 10001,
	"transaction_id": "<uuid4>"
}
```

```json
# to'lov bekor qilinganda
{
	"action":"cancel",
	"amount": 10001,
	"transaction_id": "<uuid4>"
}
```

Ikkala actionda ham callback url `200 status code` qaytarishi kerak json example

```json
{
	"ok": true
}
```

### Qo’shimcha malumot

- callback urldan success javob kelmasa dastur `3 marotaba` qayta urinadi va baribir javob success bo’lmasa transactionni yopadi.

# Muhum malumotlar

- To’lovdan avval transaction yaratasiz va dastur qaytargan miqdorda to’lov qilishini so’raysiz
- Transaction 30 daqiqa active qoladi keyin bekor qilinadi 30 daqiqadan keyingi to’langan to’lovlar tasdiqlanmaydi.
- Dastur ko’plab transactionlar bilan ishlay oladi lekin to’lov summasi farqi kattalshib ketishi mumkun masalan `10 ming` so’mlik `1000 ta` transactiondan keyin to’lov `11 ming` bo’lib ketadi buni oldini olish uchun bir nechta kartalardan foydalanishingiz mumkun dasturni bir nechta varintlarini turli accountlarga ulaysiz. (`buni hozirda qo’lda so’zlashingiz kerak  keyingi yangilanishlarda buni avtomatlashtiramiz`)

# Savollar?

- `Telegram accountni dasturga ulash hafsizmi?:` Albatta bu hafsiz chunki dastur open source ko’dlarini istalgan odam tekshirib chiqishi mumkun va malumotlar o’z serveringizda qoladi.
- `Men ham dasturni rivojlantirishga yordam bera olamanmi?:` Albatta biz yordamingizdan doim hursand bo’lamiz dasturni fork qilib oling va yangilanishlarni pull request yaratishingiz mumkun biz albatta ko’rib chiqmiz.
