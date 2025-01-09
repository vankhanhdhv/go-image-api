Đúng vậy, lỗi xảy ra là do kiểu dữ liệu của bạn khi truyền vào body trong n8n bị sai. Hiện tại, các giá trị như `title_line_height`, `content_line_height`, `title_rect`, và `content_rect` đang được truyền vào dưới dạng chuỗi (`string`), nhưng API của bạn mong đợi các kiểu dữ liệu khác, như `float64` cho line height và `[]int` cho vùng chèn (`rect`).

---

### **Hướng dẫn sửa lỗi**

#### **1. Kiểm tra định dạng dữ liệu**

Dưới đây là kiểu dữ liệu mà API của bạn mong đợi cho các tham số:

| Tham số                | Kiểu dữ liệu cần thiết       | Ví dụ                             |
|------------------------|-----------------------------|-----------------------------------|
| `title_line_height`    | `float64` (số thực)         | `1.8`                             |
| `content_line_height`  | `float64` (số thực)         | `1.6`                             |
| `title_rect`           | `[4]int` (mảng số nguyên)   | `[100, 50, 500, 150]`             |
| `content_rect`         | `[4]int` (mảng số nguyên)   | `[100, 200, 500, 300]`            |

---

#### **2. Cập nhật cấu hình n8n**

Trong n8n, bạn cần đảm bảo rằng các giá trị truyền vào body không phải là chuỗi nếu API mong đợi kiểu dữ liệu khác.

**Cách cấu hình trong n8n:**

1. **Kiểm tra giá trị truyền vào từng tham số**:
   - `title_line_height`: Đảm bảo là số thực, ví dụ: `1.8`.
   - `content_line_height`: Đảm bảo là số thực, ví dụ: `1.6`.
   - `title_rect` và `content_rect`: Đảm bảo là mảng số nguyên, ví dụ: `[100, 50, 500, 150]`.

2. **Sửa giá trị trực tiếp trong n8n:**
   - Với `title_line_height` và `content_line_height`: Loại bỏ dấu ngoặc kép (") nếu có, ví dụ: thay `"1.8"` bằng `1.8`.
   - Với `title_rect` và `content_rect`: Truyền giá trị theo dạng JSON hợp lệ. Ví dụ:
     ```
     [100, 50, 500, 150]
     ```

3. **Ví dụ cụ thể:**
   - **Title Line Height**: `1.8`
   - **Content Line Height**: `1.6`
   - **Title Rect**: `[100, 50, 500, 150]`
   - **Content Rect**: `[100, 200, 500, 300]`

---

#### **3. Nếu sử dụng n8n với biến động**

Nếu bạn đang sử dụng n8n với dữ liệu biến động (như lấy từ JSON hoặc node trước đó), hãy chắc chắn:
   - **Kiểu dữ liệu**: Đảm bảo bạn không bao quanh số thực hoặc mảng bằng dấu ngoặc kép `"`.
   - **Biểu thức động**: Sử dụng biểu thức n8n để parse giá trị đúng kiểu. Ví dụ:
     ```javascript
     {{ [100, 50, 500, 150] }}
     ```

---

### **Ví dụ Body Chuẩn**

Dưới đây là một body hợp lệ mà API của bạn mong đợi:
```json
{
  "image_url": "https://s3.amazonaws.com/i.snag.gy/BhERXe.jpg",
  "title": "What is SEO?",
  "content": "Introduction to SEO, its importance, and how it works.",
  "title_line_height": 1.8,
  "content_line_height": 1.6,
  "title_rect": [100, 50, 500, 150],
  "content_rect": [100, 200, 500, 300],
  "title_font_size": 60,
  "content_font_size": 40
}
```

---

### **Kiểm tra lại**

1. Sửa body trong n8n theo hướng dẫn trên.
2. Nhấn **Test Step** để kiểm tra.

Nếu vẫn gặp lỗi, vui lòng gửi thông tin log chi tiết hoặc cấu hình n8n để tôi hỗ trợ thêm. 😊
