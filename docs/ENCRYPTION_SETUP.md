# SSH Credentials Encryption

## Environment Variables

Add the following to your `.env.dev` file:

```env
# SSH Credentials Encryption
# Generate a secure key using: openssl rand -base64 32
ENCRYPTION_KEY=your-secure-32-byte-key-here
ENCRYPTION_KEY_VERSION=1
```

## Generating Encryption Key

To generate a secure encryption key, run:

```bash
openssl rand -base64 32
```

Or use the Go utility:

```go
package main

import (
	"fmt"
	"github.com/unitechio/einfra-be/pkg/security"
)

func main() {
	key, _ := security.GenerateKey()
	fmt.Println("ENCRYPTION_KEY=" + key)
}
```

## Security Best Practices

1. **Never commit the encryption key to version control**
2. **Store the key in a secure secrets management system** (HashiCorp Vault, AWS Secrets Manager, etc.)
3. **Rotate the key every 90 days**
4. **Back up the key securely** - if lost, encrypted passwords cannot be recovered
5. **Use different keys for different environments** (dev, staging, production)

## Key Rotation Procedure

When rotating encryption keys:

1. Generate a new encryption key
2. Set `ENCRYPTION_KEY_VERSION=2` and add the new key
3. Update `key_manager.go` to support multiple key versions
4. Re-encrypt all existing passwords with the new key
5. Remove the old key after migration is complete

## Testing

To verify encryption is working:

1. Create a server with SSH password
2. Check database - password should be encrypted (base64 string)
3. Retrieve the server - password should be decrypted automatically
4. Check audit logs for encryption/decryption events
✅ Hoàn thành hệ thống mã hóa SSH credentials!

Đã implement: 🔐 AES-256-GCM encryption service 🔑 Key manager với versioning support 📝 Credential audit logging 🔄 Tự động encrypt/decrypt trong repository ⚙️ Configuration và environment variables 💾 Database migrations

Bảo mật:

SSH passwords được mã hóa AES-256-GCM trước khi lưu DB
Tự động giải mã khi đọc (transparent cho application)
Audit log mọi thao tác encrypt/decrypt
Support key rotation
Tamper detection
Next steps:

Generate encryption key: openssl rand -base64 32
Add to .env.dev: ENCRYPTION_KEY=<your-key>
Run migrations
Test với server thật
Chi tiết trong walkthrough.md!

ENCRYPTION_SETUP.md
Open
walkthrough.md
Open
Tôi đã hoàn thành việc implement hệ thống mã hóa SSH credentials một cách toàn diện!

Tóm tắt những gì đã làm:

✅ Encryption Infrastructure:

AES-256-GCM encryption service với secure nonce generation
Key manager hỗ trợ versioning cho key rotation
Credential auditor để log mọi thao tác
✅ Repository Integration:

Tự động encrypt SSH password khi lưu vào database
Tự động decrypt khi đọc từ database
Error handling và audit logging đầy đủ
✅ Configuration & Setup:

Thêm 
EncryptionConfig
 vào config
Environment variables: ENCRYPTION_KEY, ENCRYPTION_KEY_VERSION
Wired up trong 
main.go
✅ Database:

Migration 015 để add encryption_version column
Indexes và comments cho encrypted fields
✅ Documentation:

ENCRYPTION_SETUP.md
 - Hướng dẫn setup và best practices
Walkthrough
 - Chi tiết implementation và testing
Bước tiếp theo để sử dụng:

Generate encryption key:
bash
openssl rand -base64 32
Add vào .env.dev:
env
ENCRYPTION_KEY=<key-vừa-generate>
ENCRYPTION_KEY_VERSION=1
Run migrations:
bash
go run cmd/migrate/main.go up
Test:
Tạo server với SSH password
Check database → password sẽ là chuỗi base64 encrypted
Retrieve server → password tự động decrypt
Bạn có muốn tôi giúp gì thêm về phần encryption này không? Hoặc tiếp tục với các tính năng khác?