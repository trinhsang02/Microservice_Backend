# Hướng dẫn kết nối và theo dõi database qua pgAdmin

Dự án này đã được cấu hình để sử dụng pgAdmin làm công cụ quản lý PostgreSQL qua giao diện web, giúp bạn dễ dàng theo dõi và quản lý dữ liệu.

## Khởi động các container với Docker Compose

```bash
docker-compose up -d
```

## Truy cập pgAdmin

1. Truy cập: http://localhost:5050
2. Đăng nhập với thông tin sau:
   - Email: admin@admin.com
   - Password: admin

## Cấu hình kết nối đến PostgreSQL (Lần đầu nếu truy cập không có Database)

1. Trong giao diện pgAdmin, nhấp chuột phải vào "Servers" trong cây bên trái và chọn "Create" > "Server..."
2. Trong tab "General", điền tên server (ví dụ: "Cyclone Database")
3. Chuyển sang tab "Connection" và điền thông tin sau:
   - Host name/address: postgres (tên service của PostgreSQL trong docker-compose)
   - Port: 5432
   - Maintenance database: microservices
   - Username: postgres
   - Password: postgres
4. Nhấn "Save" để lưu kết nối
