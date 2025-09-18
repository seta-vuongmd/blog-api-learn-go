# Blog API - Hiệu năng cao

## Mô tả
API Blog sử dụng Go, PostgreSQL, Redis, Elasticsearch. Đáp ứng các yêu cầu về hiệu năng, caching, full-text search.

## Khởi chạy hệ thống

```sh
docker-compose up --build
```

## Các endpoint chính

- Đăng ký: `POST /register`
- Đăng nhập: `POST /login`
- Tạo bài viết: `POST /posts`
- Lấy danh sách bài viết: `GET /posts`
- Lấy chi tiết bài viết: `GET /posts/:id`
- Tìm kiếm theo tag: `GET /posts/search-by-tag?tag=<tag_name>`
- Tìm kiếm full-text: `GET /posts/search?q=<query_string>`
- Cập nhật bài viết: `PUT /posts/:id`
- Xóa bài viết: `DELETE /posts/:id`

## Ví dụ curl

```sh
# Đăng ký
curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d '{"username":"user1","password":"pass"}'

# Đăng nhập
curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"username":"user1","password":"pass"}'

# Tạo bài viết
curl -X POST http://localhost:8080/posts -H "Content-Type: application/json" -d '{"title":"Tiêu đề","content":"Nội dung","tags":["go","api"],"user_id":1}'

# Lấy danh sách bài viết
curl http://localhost:8080/posts

# Lấy chi tiết bài viết
curl http://localhost:8080/posts/1

# Cập nhật bài viết
curl -X PUT http://localhost:8080/posts/1 -H "Content-Type: application/json" -d '{"title":"Tiêu đề mới","content":"Nội dung mới","tags":["go","update"],"user_id":1}'

# Xóa bài viết
curl -X DELETE http://localhost:8080/posts/1

# Tìm kiếm theo tag
curl http://localhost:8080/posts/search-by-tag?tag=go

# Tìm kiếm full-text
curl http://localhost:8080/posts/search?q=Tiêu đề

# Lấy bài viết liên quan
curl http://localhost:8080/posts/1/related
```

## Ghi chú
- Đảm bảo các service đã chạy trước khi gọi API.
- Cấu hình kết nối DB, Redis, Elasticsearch có thể chỉnh sửa trong code.
