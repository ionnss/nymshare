![[logo.png]]


**A Secure Ephemeral File-Sharing Platform for Whistleblowers and Journalists**

> *Whispers encrypted, secrets ephemeral.*

---

### **Super Simple Explanation**  
**NymShare does two things**:  
1. Lets whistleblowers **encrypt files** so **only one journalist** can open them.  
2. **Deletes those files** after 24 hours.  

---

### **Step-by-Step Workflow**  
#### **1. Journalist Setup**  
- **What**: A journalist (or media org) wants to receive files securely.  
- **How**:  
  - They generate **two keys** on their own device (like a password pair):  
    - **Public Key**: Like a "lock." Shared publicly as a QR code/URL.  
    - **Private Key**: Like a "key." Never leaves their device.  
  - They share the **public key** (QR code/URL) on their profile (e.g., social media, website).  

---

#### **2. Whistleblower Upload**  
- **What**: A whistleblower wants to send a file to that journalist.  
- **How**:  
  - They scan the **journalist’s QR code** (or visit the URL) to get the public key.  
  - The whistleblower’s browser **encrypts the file** using the public key (*like locking a box*).  
  - The encrypted file is uploaded to NymShare’s server.  

---

#### **3. Delivery & Access**  
- **What**: The journalist needs to open the encrypted file.  
- **How**:  
  - NymShare sends the journalist a ProtonMail alert: *“You have a new file.”*  
  - The journalist visits the link and **uses their private key** (stored on their device) to **decrypt the file**.  
  - After 24 hours, the file is **automatically deleted** from the server.  

---

### **Key Clarifications**  
#### **Why Keys?**  
- **Public Key** = Lock. Anyone can use it to lock (encrypt) a file.  
- **Private Key** = Key. Only the journalist can unlock (decrypt) it.  
- **No logins, no accounts**: The journalist’s "identity" is their ProtonMail (to receive alerts) and their key pair.  

---

#### **Where Keys Are Generated**  
- **Journalist’s Device**: Keys are created on their laptop/phone (not on your server).  
  - Example Tools:  
    - Browser: `window.crypto.subtle.generateKey()` (JavaScript).  
    - Apps: Libraries like `libsodium` or `OpenSSL`.  
- **Whistleblower’s Device**: Encryption happens in their browser. Your server **never sees unencrypted files**.  

---

#### **No Sessions, No Logins**  
- **Whistleblowers**: Don’t need to sign up. Just upload files.  
- **Journalists**: Only need to provide a ProtonMail address to receive alerts.  

---

### **Visual Workflow**  
```plaintext
Journalist:
1. Generates 🔑 (private) and 🔒 (public) on their device.
2. Shares 🔒 (QR code/URL) publicly.

Whistleblower:
1. Scans 🔒 (QR code/URL).
2. Encrypts file with 🔒 ➔ 🔐 (encrypted file).
3. Uploads 🔐 to NymShare.

NymShare:
1. Stores 🔐 for 24h.
2. Sends ProtonMail alert to journalist.

Journalist:
1. Accesses 🔐 via link in email.
2. Uses 🔑 to decrypt 🔐 ➔ 📄 (original file).
3. After 24h, 🔐 is destroyed. 💥
```

---

### **Technical Minimalism**  
If you want to skip complex encryption logic:  
1. Use **age** (a simple encryption tool):  
   - Journalist runs `age-keygen` to generate keys.  
   - Whistleblower encrypts files with `age -e -R public_key.txt file > file.age`.  
2. Host encrypted files on a server with a **24h auto-delete cron job**.  

---

### **What You Need to Build**  
1. **Frontend**:  
   - A webpage where whistleblowers drag-and-drop files, scan QR codes, and encrypt files (using JavaScript).  
2. **Backend**:  
   - A server to store encrypted files temporarily (no metadata).  
   - ProtonMail integration to send alerts.  

---

