# 🔐 npm 2FA Setup & Publishing

## Issue: 403 Forbidden - Two-factor authentication required

npm now requires 2FA for publishing packages. Here's how to set it up:

---

## ✅ Solution 1: Enable 2FA (Recommended)

### Step 1: Enable 2FA on your npm account

1. Go to: https://www.npmjs.com/settings/rudra_826/tfa
2. Click "Enable Two-Factor Authentication"
3. Choose your method:
   - **Authenticator App** (Google Authenticator, Authy, etc.) - Recommended
   - **SMS** (less secure)
4. Scan QR code or enter code manually
5. Save recovery codes in a safe place!

### Step 2: Publish with OTP

After enabling 2FA:

```bash
cd D:\ecrypto\npm

# Method 1: Provide OTP with command
npm publish --otp=123456
# Replace 123456 with your 6-digit code from authenticator app

# Method 2: npm will prompt you
npm publish
# Then enter OTP when asked
```

---

## ✅ Solution 2: Granular Access Token (For CI/CD)

If you need to automate publishing:

### Step 1: Create a granular token

1. Go to: https://www.npmjs.com/settings/rudra_826/tokens
2. Click "Generate New Token" → "Granular Access Token"
3. Configure:
   - **Name**: "ecrypto-cli-publish"
   - **Expiration**: 90 days
   - **Packages and scopes**: Select "ecrypto-cli"
   - **Permissions**: Read and write
4. Enable "Bypass 2FA" (only for automation)
5. Copy the token

### Step 2: Use the token

```bash
# Set token in .npmrc
npm config set //registry.npmjs.org/:_authToken YOUR_TOKEN

# Or use environment variable
export NPM_TOKEN=YOUR_TOKEN

# Then publish
npm publish
```

---

## 🚀 Quick Fix Right Now

**Easiest way:**

1. Install authenticator app on your phone:

   - Google Authenticator
   - Microsoft Authenticator
   - Authy

2. Visit: https://www.npmjs.com/settings/rudra_826/tfa

3. Enable 2FA with app

4. Get your 6-digit code from app

5. Publish:
   ```bash
   cd D:\ecrypto\npm
   npm publish --otp=YOUR_CODE
   ```

---

## 📱 Recommended Authenticator Apps

- **Google Authenticator** (Android, iOS)
- **Microsoft Authenticator** (Android, iOS)
- **Authy** (Android, iOS, Desktop)
- **1Password** (if you use it)

---

## ⚠️ Important Notes

1. **Save recovery codes**: npm will give you recovery codes - save them!
2. **Keep OTP codes handy**: You'll need them for every publish
3. **Token expiration**: Granular tokens expire after 90 days
4. **Security**: Don't share tokens or disable 2FA

---

## 🔍 Troubleshooting

### Error: "Invalid OTP"

- Make sure your phone's time is synced
- Code expires every 30 seconds, get a fresh one

### Error: "Token expired"

- Create a new granular token
- Or re-enable 2FA and use OTP

### Can't access authenticator app

- Use recovery codes
- Contact npm support

---

## ✅ After Setup

Once 2FA is enabled, your publishing command is:

```bash
cd D:\ecrypto\npm
npm publish --otp=123456
```

Replace `123456` with the code from your authenticator app.

---

## 📞 Need Help?

- npm 2FA docs: https://docs.npmjs.com/configuring-two-factor-authentication
- npm support: https://www.npmjs.com/support

---

**Ready to enable 2FA?**

Visit: https://www.npmjs.com/settings/rudra_826/tfa

Then publish with: `npm publish --otp=YOUR_CODE`
