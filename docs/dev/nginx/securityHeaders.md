# Security Headers (NGINX)

## Overview

Security headers are HTTP response headers that instruct the browser to enforce certain security behaviors.
They do not directly secure the backend, but help protect users from common web-based attacks.

In this project, these headers are configured at the NGINX layer (reverse proxy), ensuring all responses include them.

---

## Headers Used

### 1. X-Frame-Options

```
X-Frame-Options: DENY
```

**Purpose:** Prevents clickjacking attacks by disallowing the application from being embedded in an iframe.

---

### 2. X-Content-Type-Options

```
X-Content-Type-Options: nosniff
```

**Purpose:** Prevents browsers from guessing the content type (MIME sniffing).
Ensures files are treated as declared.

---

### 3. X-XSS-Protection

```
X-XSS-Protection: 1; mode=block
```

**Purpose:** Enables basic XSS filtering in older browsers.
Modern browsers rely more on Content Security Policy (CSP).

---

### 4. Strict-Transport-Security (HSTS)

```
Strict-Transport-Security: max-age=31536000; includeSubDomains
```

**Purpose:** Forces browsers to use HTTPS for all future requests.
Prevents downgrade attacks (HTTP instead of HTTPS).

⚠️ Should only be enabled when HTTPS is correctly configured.

---

### 5. Referrer-Policy

```
Referrer-Policy: no-referrer-when-downgrade
```

**Purpose:** Controls how much referrer information is sent when navigating between pages.
Prevents leaking sensitive URLs to less secure destinations.

---

## Why Configure at NGINX?

- Centralized enforcement for all services
- Keeps application code clean
- Applies to all responses (API + static content)
- Matches real-world production architecture

---

## Notes

- These headers protect the **client/browser**, not the backend directly
- They are part of **defense in depth**
- More advanced setups may include:
  - Content-Security-Policy (CSP)
  - Permissions-Policy
  - Rate limiting at proxy level

---

## Verification

Run:

```bash
curl -I https://localhost/tasks -k
```

Expected output should include:

- X-Frame-Options
- X-Content-Type-Options
- Strict-Transport-Security
- Referrer-Policy

---

## Takeaway

Security headers are a low-cost, high-impact way to improve application security
and are standard in production environments.
