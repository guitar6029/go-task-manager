## ⚠️ Windows / Git Bash Notes (OpenSSL Path Issue)

When generating SSL certificates using `openssl` in Git Bash (MSYS/Cygwin), paths that start with `/` may be incorrectly converted into Windows-style paths.

### ❌ Problem

Running:

```
openssl req ... -subj "/CN=localhost"
```

may result in an error like:

```
subject name is not in that format: 'C:/Program Files/Git/CN=localhost'
```

### ✅ Solution

Disable path conversion for this command:

```
MSYS_NO_PATHCONV=1 openssl req -x509 -nodes -days 365 \
  -newkey rsa:2048 \
  -keyout nginx/certs/key.pem \
  -out nginx/certs/cert.pem \
  -subj "/CN=localhost"
```

### 🧠 Explanation

Git Bash automatically rewrites paths like `/something` into Windows paths (e.g., `C:\...`).
This breaks tools like OpenSSL that expect Unix-style input.

Setting `MSYS_NO_PATHCONV=1` prevents this behavior.

### 🔁 Alternative

Use WSL (Windows Subsystem for Linux) to avoid these issues entirely when working with Unix-based tools.

---
