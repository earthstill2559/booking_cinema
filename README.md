# Cinema Booking System

ระบบจองที่นั่งโรงภาพยนตร์แบบ Real-time พร้อม Redis Distributed Lock, WebSocket, MongoDB Audit Logs และ RabbitMQ

## Features

### User Side
- **Authentication**: Google OAuth / Firebase (หรือ Dev Login ใน DEV_MODE)
- **Real-time Seat Map**: สถานะ `AVAILABLE`, `LOCKED`, `BOOKED` ผ่าน WebSocket
- **Booking Flow**: Lock ที่นั่ง 5 นาที → ชำระเงิน → BOOKED / หมดเวลา → ปล่อย lock

### Admin Side
- Dashboard ดูรายการจองทั้งหมด
- Filter ตาม movie, date, user_id
- API ป้องกันด้วย role `ADMIN`

### Infrastructure
- **Redis**: Distributed Lock (SETNX + TTL 5 นาที)
- **MongoDB**: Bookings + Audit Logs
- **RabbitMQ**: ส่ง event เมื่อจองสำเร็จ → async notification (mock) + audit log

## Lock Strategy

เราใช้ **Redis SETNX + TTL** เป็น Distributed Lock:

| รายการ | รายละเอียด |
|--------|------------|
| Key | `lock:{showtime_id}:{seat_id}` |
| Value | `user_id` ของผู้ lock |
| TTL | 5 นาที |
| Acquire | `SETNX` — atomic, มีเพียง 1 client ที่ lock ได้ |
| Release | `DEL` เมื่อยกเลิกหรือจองสำเร็จ |
| Expiry | Redis keyspace notification → broadcast `AVAILABLE` + audit log |

**เหตุผลที่เลือก SETNX + TTL**
1. **Atomic** — ป้องกัน double booking เมื่อหลาย user กดพร้อมกัน
2. **Self-healing** — TTL ปล่อย lock อัตโนมัติถ้าไม่ชำระเงินภายใน 5 นาที
3. **Distributed** — ใช้ได้กับหลาย instance ของ backend
4. **Owner tracking** — เก็บ user_id ใน value เพื่อ verify ก่อน confirm

**Concurrency**: MongoDB unique index บน `(showtime_id, seat_id)` สำหรับ status CONFIRMED เป็น safety net ชั้นที่ 2

## Audit Logs (MongoDB)

| Event | เมื่อไหร่ |
|-------|-----------|
| `BOOKING_SUCCESS` | จองสำเร็จ |
| `BOOKING_TIMEOUT` | Lock หมดเวลา 5 นาที |
| `SEAT_RELEASED` | User ยกเลิก หรือ lock หมดอายุ |
| `SYSTEM_ERROR` | Lock fail / Redis error |

## Message Queue (RabbitMQ)

Use case: เมื่อจองสำเร็จ → publish `booking.success` → consumer ทำ async notification (mock email) + audit log

## Run with Docker

```bash
docker compose up --build
```

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- RabbitMQ Management: http://localhost:15672 (guest/guest)

### Dev Login (DEV_MODE=true)
- User: `user001` — จองที่นั่ง
- Admin: `admin001` — เข้า Admin Dashboard

## Local Development

### Backend
```bash
cd backend
cp .env.example .env
# Start Redis, MongoDB, RabbitMQ (via docker compose up redis mongo rabbitmq)
go run .
```

### Frontend
```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/showtime` | - | ข้อมูลรอบฉาย |
| GET | `/api/seats` | - | แผนที่ที่นั่ง |
| GET | `/ws` | - | WebSocket real-time |
| GET | `/api/me` | User | ข้อมูล user + role |
| POST | `/api/seats/:id/lock` | User | Lock ที่นั่ง |
| POST | `/api/seats/:id/unlock` | User | ปล่อย lock |
| POST | `/api/bookings/confirm` | User | ยืนยันการจอง (payment) |
| GET | `/api/admin/bookings` | Admin | รายการจอง + filter |
| GET | `/api/admin/audit-logs` | Admin | Audit logs |

## Environment Variables

ดู `.env.example` ใน `backend/` และ `frontend/`

## Security

- Role separation: `USER` / `ADMIN`
- Admin APIs ใช้ middleware `AuthRequired` + `AdminRequired`
- Firebase ID token verification (production) หรือ dev token (development)
