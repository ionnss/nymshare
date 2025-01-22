![[logo.png]]


**A Secure Ephemeral File-Sharing Platform for Whistleblowers and Journalists**

> *Whispers encrypted, secrets ephemeral.*

---

Here’s the **revised technical documentation** using **Go (Golang)** for the backend and minimal JavaScript only where absolutely necessary (for client-side decryption). Decryption must happen in the browser, so JavaScript is unavoidable there, but we’ll minimize its use elsewhere:

---

# **NymShare Technical Design**  
*Secure Ephemeral File Sharing with Go, HTML, and JavaScript*

---

## **1. Core Workflow**  
### **1.1 Journalist Registration**  
- **Journalist Input**: Provides ProtonMail address.  
- **Key Generation**: Done client-side (JavaScript) to avoid server handling private keys.  
- **Database (Go + SQLite)**:  
  ```go
  type Journalist struct {
    ID         int    `gorm:"primaryKey"`
    ProtonMail string `gorm:"unique"`
    PublicKey  string // Base64-encoded public key
  }
  ```

---

## **2. Technical Stack**  
### **2.1 Backend (Go)**  
- **HTTP Server**: `net/http` or `Gin` for routing.  
- **Database**: SQLite with `gorm.io/gorm`.  
- **File Storage**: Local disk or S3-compatible storage.  
- **Email**: ProtonMail SMTP integration.  

### **2.2 Frontend**  
- **HTML**: Static templates for registration, upload, and decryption.  
- **JavaScript**: Minimal code for key generation, encryption, and decryption (using WebCrypto API).  

---

## **3. Component Breakdown**  
### **3.1 Journalist Registration (Go + JS)**  
#### **Go Handler**  
```go
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    protonMail := r.FormValue("protonmail")
    publicKey := r.FormValue("public_key") // Sent from client-side JS

    // Store in SQLite
    db.Create(&Journalist{ProtonMail: protonMail, PublicKey: publicKey})
    
    fmt.Fprintf(w, "Registration successful. Download your private key.")
}
```

#### **Client-Side JavaScript (Key Generation)**  
```html
<script>
async function generateKeys() {
    // Generate RSA-OAEP key pair
    const keyPair = await crypto.subtle.generateKey(
        { name: "RSA-OAEP", modulusLength: 4096, publicExponent: new Uint8Array([1, 0, 1]), hash: "SHA-256" },
        true, ["encrypt", "decrypt"]
    );

    // Export public key (to send to server)
    const publicKey = await crypto.subtle.exportKey("spki", keyPair.publicKey);
    const publicKeyBase64 = btoa(String.fromCharCode(...new Uint8Array(publicKey)));
    document.getElementById("publicKey").value = publicKeyBase64;

    // Export private key (for download)
    const privateKey = await crypto.subtle.exportKey("pkcs8", keyPair.privateKey);
    const privateKeyBlob = new Blob([privateKey], { type: "application/octet-stream" });
    const url = URL.createObjectURL(privateKeyBlob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "private_key.pem";
    a.click();
}
</script>
```

---

### **3.2 Whistleblower Upload (Go)**  
#### **Go File Upload Handler**  
```go
func UploadHandler(w http.ResponseWriter, r *http.Request) {
    // Fetch journalist's public key from DB
    journalistID := r.URL.Query().Get("id")
    var journalist Journalist
    db.First(&journalist, journalistID)

    // Save encrypted file
    file, header, _ := r.FormFile("file")
    defer file.Close()
    encryptedFilePath := fmt.Sprintf("./uploads/%s_%s", journalistID, header.Filename)
    out, _ := os.Create(encryptedFilePath)
    defer out.Close()
    io.Copy(out, file)

    // Schedule deletion in 24h
    time.AfterFunc(24*time.Hour, func() { os.Remove(encryptedFilePath) })

    // Send ProtonMail alert
    sendProtonMailAlert(journalist.ProtonMail, encryptedFilePath)
}
```

---

### **3.3 Decryption & Download (JavaScript)**  
#### **Decryption Page (HTML/JS)**  
```html
<!DOCTYPE html>
<html>
<body>
    <input type="file" id="privateKeyFile">
    <button onclick="decryptAndDownload()">Decrypt</button>

    <script>
    async function decryptAndDownload() {
        // Fetch encrypted file from Go server
        const fileID = window.location.pathname.split("/").pop();
        const encryptedFile = await fetch(`/files/${fileID}`);
        const encryptedData = await encryptedFile.arrayBuffer();

        // Read private key
        const privateKeyFile = document.getElementById("privateKeyFile").files[0];
        const privateKeyBuffer = await privateKeyFile.arrayBuffer();
        
        // Import private key
        const privateKey = await crypto.subtle.importKey(
            "pkcs8",
            privateKeyBuffer,
            { name: "RSA-OAEP", hash: "SHA-256" },
            true,
            ["decrypt"]
        );

        // Decrypt
        const decryptedData = await crypto.subtle.decrypt(
            { name: "RSA-OAEP" },
            privateKey,
            encryptedData
        );

        // Trigger download
        const blob = new Blob([decryptedData]);
        const url = URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = "decrypted_file";
        a.click();
    }
    </script>
</body>
</html>
```

---

## **4. Security Design**  
### **4.1 Threat Mitigation**  
- **No Private Key on Server**: Journalists’ private keys never leave their devices.  
- **Ephemeral Storage**: Files auto-delete after 24h via Go’s `time.AfterFunc`.  
- **TLS/HTTPS**: Enforced for all traffic (use Let’s Encrypt).  

### **4.2 ProtonMail Integration (Go)**  
```go
func sendProtonMailAlert(to, fileID string) {
    auth := smtp.PlainAuth("", "your@protonmail.com", "password", "smtp.protonmail.com")
    msg := fmt.Sprintf("From: NymShare <noreply@nymshare.com>\nTo: %s\nSubject: New Secure File\n\nDownload: https://nymshare.com/files/%s", to, fileID)
    
    err := smtp.SendMail("smtp.protonmail.com:587", auth, "noreply@nymshare.com", []string{to}, []byte(msg))
    if err != nil {
        log.Printf("Failed to send email: %v", err)
    }
}
```

---

## **5. Deployment**  
### **5.1 Go Server Setup**  
```bash
go mod init nymshare
go get github.com/gin-gonic/gin gorm.io/gorm gorm.io/driver/sqlite
go run main.go
```

### **5.2 Frontend**  
- Serve static HTML/JS files via Go’s `http.FileServer`.  

---

## **6. Why JavaScript is Unavoidable**  
- **Browser Limitations**: Go cannot run in the browser (outside WebAssembly, which is overly complex here).  
- **WebCrypto API**: Required for client-side encryption/decryption. JavaScript is the only way to access it.  

---

## **7. Alternatives to JavaScript**  
If avoiding JavaScript is critical:  
1. **Desktop/Mobile App**: Build a Go app (using Fyne or Wails) to handle encryption/decryption.  
2. **WebAssembly (WASM)**: Compile Go to WASM for browser execution (complex setup).  

---

This design prioritizes **security** and **simplicity**, using Go for the backend and minimal JavaScript for client-side cryptography. Let me know if you’d like to refine any component! 🔐