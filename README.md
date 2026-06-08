# URL Shortening Service

## คำอธิบาย Project

Project นี้เป็นบริการแปลง URL ยาวเป็น URL สั้นๆ (URL Shortening Service) ที่พัฒนาด้วย Go และ MongoDB

### ฟีเจอร์หลัก:
- **สร้าง URL สั้น**: แปลง URL ยาวเป็น URL สั้นพร้อม ID เฉพาะ
- **เปลี่ยนเส้นทาง**: ใช้ URL สั้นเพื่อเปลี่ยนเส้นทางไปยัง URL ต้นฉบับ
- **ติดตามการเข้าถึง**: นับจำนวนครั้งที่แต่ละ URL ได้รับการเข้าถึง
- **ลบ URL**: ลบ URL สั้นที่ไม่ต้องการใช้แล้ว
- **แยกชั้น logic**: ใช้ `services` สำหรับ business logic, `handlers` สำหรับ request/response เท่านั้น

### เทคโนโลยีที่ใช้:
- **Go 1.26.3**: ภาษาโปรแกรมมิ่ง
- **Gin Web Framework**: Framework สำหรับสร้าง REST API
- **MongoDB**: ฐานข้อมูล NoSQL
- **KSUID**: สำหรับสร้าง ID เฉพาะ
- **Docker**: สำหรับการ Deploy และ Containerization

---

## วิธี Run Project

### ข้อกำหนดเบื้องต้น:
- Go 1.26.3 ขึ้นไป
- Docker และ Docker Compose
- MongoDB (หรือใช้ Docker สำหรับรัน MongoDB)

### ขั้นตอนการรัน:

#### **วิธีที่ 1: ใช้ Docker Compose (ง่ายที่สุด)**

1. **ไปที่โฟลเดอร์ backend:**
   ```bash
   cd backend
   ```

2. **รัน Docker Compose:**
   ```bash
   docker compose -p shorturl-service up --build -d
   ```
   
   คำสั่งนี้จะ:
   - เปิดตัว MongoDB Container
   - สร้าง Build และรัน URL Shortening Service

3. **เมื่อเห็นข้อความ "Server starting on port 8080" แสดงว่าเซิร์ฟเวอร์ทำงานแล้ว**

#### **วิธีที่ 2: รัน Go Application โดยตรง**

1. **ติดตั้ง Go Dependencies:**
   ```bash
   cd backend
   go mod download && go mod verify
   ```

2. **สร้างไฟล์ .env (ตั้งค่าสิ่งแวดล้อม):**
   ```bash
   cat > .env << EOF
   MONGODB_URI=mongodb://localhost:27017
   DB_NAME=url_shortening_db
   PORT=8080
   EOF
   ```

3. **เปิดตัว MongoDB:**
   ```bash
   docker run -d -p 27017:27017 --name mongodb mongo:latest
   ```

4. **รัน Go Application:**
   ```bash
   go run main.go
   ```

---

## API Endpoints

| Method | Endpoint | คำอธิบาย |
|--------|----------|---------|
| POST | `/api/v1/shortens` | สร้าง URL สั้น |
| GET | `/api/v1/shortens/:shortCode` | เปลี่ยนเส้นทางไปยัง URL ต้นฉบับ และติดตามการเข้าถึง |
| PUT | `/api/v1/shortens/:shortCode` | อัปเดต URL ต้นฉบับ |
| DELETE | `/api/v1/shortens/:shortCode` | ลบ URL สั้น |
| GET | `/api/v1/shortens/:shortCode/stats` | ดูสถิติการเข้าถึง URL |

### ตัวอย่างการใช้:

**สร้าง URL สั้น:**
```bash
curl -X POST http://localhost:8080/api/v1/shortens \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.example.com/very/long/url/path"
  }'
```

**เปิด URL สั้น (เปลี่ยนเส้นทาง):**
```bash
curl -X GET http://localhost:8080/api/v1/shortens/abc123
```

**อัปเดต URL ต้นฉบับ:**
```bash
curl -X PUT http://localhost:8080/api/v1/shortens/abc123 \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.example.com/new/url"
  }'
```

**ดูสถิติการเข้าถึง:**
```bash
curl -X GET http://localhost:8080/api/v1/shortens/abc123/stats
```

**ลบ URL สั้น:**
```bash
curl -X DELETE http://localhost:8080/api/v1/shortens/abc123
```

---

## โครงสร้าง Project

```
backend/
├── main.go                    # Entry point ของ Application
├── go.mod                     # Go Module dependencies
├── docker-compose.yml         # Docker Compose configuration
├── config/
│   └── config.go             # ตั้งค่า MongoDB connection
├── handlers/
│   └── short_url_handler.go  # HTTP request handlers
├── services/
│   └── service.go            # Business logic layer
├── models/
│   └── short_url.go          # Data models
├── repository/
│   └── short_url_repository.go # Database access layer
└── routes/
    └── routes.go             # API routes configuration
```

---

## ตัวแปร Environment

สามารถตั้งค่าผ่าน `.env` file:

| ตัวแปร | ค่าเริ่มต้น | คำอธิบาย |
|--------|-----------|---------|
| `MONGODB_URI` | `mongodb://localhost:27017` | Connection string ของ MongoDB |
| `DB_NAME` | `gin_mongodb_api` | ชื่อฐานข้อมูล |
| `PORT` | `8080` | Port ของ API Server |

---

## การหยุด Service

**ถ้าใช้ Docker Compose:**
```bash
docker-compose down
```

**ถ้ารัน Go Application โดยตรง:**
- กด `Ctrl + C` ในเทอร์มินัล

---

## ไฟล์ที่สำคัญ

- **main.go**: กำหนดค่าเซิร์ฟเวอร์ Gin, เชื่อมต่อ MongoDB, และเริ่มต้น handlers
- **config.go**: จัดการการเชื่อมต่อ MongoDB
- **short_url_handler.go**: ประมวลผล HTTP requests สำหรับสร้าง ดู และลบ URLs
- **service.go**: เก็บ business logic สำหรับการสร้าง, อัปเดต, ลบ, และเรียกดูข้อมูล short URL
- **short_url_repository.go**: ติดต่อกับ MongoDB สำหรับ CRUD operations
- **routes.go**: กำหนด API endpoints

---