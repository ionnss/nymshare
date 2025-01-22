![[logo.png]]


**A Secure Ephemeral File-Sharing Platform for Whistleblowers and Journalists**

> *Whispers encrypted, secrets ephemeral.*

---

Here’s a **rewrite of the workflow and technical design** using **Go (Golang)**, **HTML templates**, and **HTMX** for a minimal, secure, and modern implementation:

---

# **NymShare Technical Design (Go + HTMX)**  
*Secure Ephemeral File Sharing for Whistleblowers & Journalists*

---

## **1. Journalist Registration**  
### **Go Backend (Handler)**  
```go
// Journalist struct (SQLite model)
type Journalist struct {
    ID         int
    ProtonMail string `gorm:"unique"`
    PublicKey  string
}

// Register journalist (HTMX-driven form)
func RegisterJournalist(c *gin.Context) {
    protonMail := c.PostForm("protonmail")
    publicKey := c.PostForm("public_key") // Client generates key pair
    
    // Store in SQLite
    db.Create(&Journalist{ProtonMail: protonMail, PublicKey: publicKey})
    
    // HTMX response: Show success message
    c.Header("HX-Trigger-After-Settle", "registrationSuccess")
    c.String(http.StatusOK, "Registration successful! Download your private key.")
}
```

### **HTML/HTMX Registration Form**  
```html
<!-- templates/register.html -->
<form hx-post="/register" hx-target="#registration-result">
    <input type="email" name="protonmail" placeholder="ProtonMail" required>
    <div id="key-generation">
        <button type="button" onclick="generateKeyPair()">Generate Keys</button>
        <input type="hidden" name="public_key" id="publicKey">
    </div>
    <div id="registration-result"></div>
</form>

<script>
// Client-side key generation (WebCrypto)
async function generateKeyPair() {
    const keyPair = await window.crypto.subtle.generateKey(
        { name: "RSA-OAEP", modulusLength: 4096, publicExponent: new Uint8Array([1,0,1]), hash: "SHA-256" },
        true, ["encrypt", "decrypt"]
    );
    
    // Export public key and set as hidden input
    const publicKey = await window.crypto.subtle.exportKey("spki", keyPair.publicKey);
    document.getElementById("publicKey").value = btoa(String.fromCharCode(...new Uint8Array(publicKey)));
    
    // Trigger private key download
    const privateKey = await window.crypto.subtle.exportKey("pkcs8", keyPair.privateKey);
    downloadPrivateKey(privateKey);
}
</script>
```

---

## **2. Whistleblower Upload**  
### **Go Backend (API)**  
```go
// Upload encrypted file
func UploadFile(c *gin.Context) {
    file, _ := c.FormFile("encryptedFile")
    journalistID := c.Param("id")
    
    // Store encrypted file (e.g., local disk or S3)
    path := fmt.Sprintf("./uploads/%s_%s", journalistID, file.Filename)
    c.SaveUploadedFile(file, path)
    
    // Schedule deletion after 24h
    go scheduleDeletion(path, 24*time.Hour)
    
    // Send ProtonMail alert
    var journalist Journalist
    db.First(&journalist, journalistID)
    sendProtonMailAlert(journalist.ProtonMail)
}
```

### **HTML/HTMX Upload Page**  
```html
<!-- templates/upload.html -->
<div hx-target="this" hx-swap="outerHTML">
    <div hx-get="/journalist/{{.JournalistID}}" hx-trigger="load"></div>
    
    <form hx-encoding="multipart/form-data" hx-post="/upload/{{.JournalistID}}">
        <input type="file" name="encryptedFile" id="fileInput" 
               onchange="encryptAndUpload(this.files[0])">
    </form>
</div>

<script>
// Client-side encryption (WebCrypto + HTMX)
async function encryptAndUpload(file) {
    const publicKey = await loadPublicKey(); // From fetched HTML
    const encryptedFile = await encryptFile(file, publicKey);
    
    // HTMX upload
    const formData = new FormData();
    formData.append("encryptedFile", encryptedFile);
    htmx.ajax("POST", "/upload/{{.JournalistID}}", { body: formData });
}
</script>
```

---

## **3. Journalist Decryption & Download**  
### **Go Backend (Serve Static Page)**  
```go
func DecryptPage(c *gin.Context) {
    c.HTML(http.StatusOK, "decrypt.html", gin.H{
        "FileID": c.Param("id"),
    })
}
```

### **HTML/JavaScript Decryption Page**  
```html
<!-- templates/decrypt.html -->
<input type="file" id="privateKey" accept=".pem">
<button onclick="decryptFile('{{.FileID}}')">Decrypt & Download</button>

<script>
async function decryptFile(fileID) {
    const privateKeyFile = document.getElementById("privateKey").files[0];
    const privateKey = await privateKeyFile.text();
    
    // Fetch encrypted file from Go server
    const encryptedFile = await fetch(`/files/${fileID}`);
    
    // Client-side decryption
    const decrypted = await window.crypto.subtle.decrypt(
        { name: "RSA-OAEP" },
        privateKey,
        await encryptedFile.arrayBuffer()
    );
    
    // Trigger download
    const blob = new Blob([decrypted]);
    window.location.href = URL.createObjectURL(blob);
}
</script>
```

---

## **4. Security & HTMX Enhancements**  
### **HTMX Security Headers (Go Middleware)**  
```go
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Content-Security-Policy", "default-src 'self'")
        c.Header("X-Frame-Options", "DENY")
        c.Next()
    }
}
```

### **Ephemeral File Cleanup (Go Cron)**  
```go
func scheduleDeletion(path string, delay time.Duration) {
    time.AfterFunc(delay, func() {
        os.Remove(path)
    })
}
```

---

## **5. ProtonMail Integration (Go SMTP)**  
```go
func sendProtonMailAlert(to string) {
    auth := smtp.PlainAuth("", "your@protonmail.com", "password", "smtp.protonmail.com")
    msg := fmt.Sprintf("Subject: New Secure File\n\nVisit: https://nymshare.com/files/%s", fileID)
    
    err := smtp.SendMail("smtp.protonmail.com:587", auth, "noreply@nymshare.com", []string{to}, []byte(msg))
    if err != nil {
        log.Printf("Failed to send email: %v", err)
    }
}
```

---

## **Key Advantages of This Stack**  
1. **Performance**: Go handles encryption/decryption I/O efficiently.  
2. **Simplicity**: HTMX eliminates complex JavaScript frameworks.  
3. **Security**:  
   - Zero private key exposure (client-side decryption).  
   - Encrypted files never touch disk unencrypted.  
4. **Portability**: Single binary deployment with SQLite.  

---

## **To Run**  
```bash
go mod init nymshare
go get github.com/gin-gonic/gin gorm.io/gorm htmx.org
go run main.go
```

This design maintains **end-to-end encryption** while leveraging Go’s concurrency and HTMX’s lightweight interactivity. Would you like me to expand on any component? 🚀