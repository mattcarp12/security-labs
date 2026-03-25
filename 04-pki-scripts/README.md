# Mini-PKI: A Python-Based Certificate Authority Lab

This repository contains a suite of Python scripts designed to demystify **Public Key Infrastructure (PKI)**. By following these labs, you will build a functional "Chain of Trust" from scratch—moving from raw mathematical keys to a signed X.509 certificate.

## 🚀 Lab Roadmap

### 1. Key Generation (`rsa.py` & `ecdsa.py`)
Generate the raw mathematical material needed for encryption and signing.
* **RSA:** Based on the difficulty of factoring large prime numbers.
* **ECDSA:** Based on the algebraic structure of elliptic curves.
* **How to run:** 
    ```bash
    python rsa.py     # Select 2048, 3072, or 4096 bits
    python ecdsa.py   # Select SECP256R1, 384, or 521
    ```

### 2. The Root CA (`ca_root.py`)
Transform your private key into a **Trust Anchor**. This script creates a self-signed certificate that acts as your personal "Identity Provider."
* **How to run:** `python ca_root.py`
* **Output:** `root_cert.pem`

### 3. The Application (`generate_csr.py`)
Simulate a server (like a web server or a Kubernetes pod) requesting an identity. This creates a **Certificate Signing Request (CSR)**.
* **How to run:** `python generate_csr.py`
* **Output:** `server_request.csr`

### 4. The Issuance (`sign_csr.py`)
Play the role of the CA. You use the Root CA's private key to sign the CSR, issuing a final, valid certificate.
* **How to run:** `python sign_csr.py`
* **Output:** `server_signed.crt`

---

## 🧠 How the Cryptography Works

### 1. RSA vs. ECDSA
* **RSA (Rivest-Shamir-Adleman):** Security relies on the product of two massive primes ($n = p \times q$). It is the "gold standard" for compatibility but requires large keys (e.g., 2048 bits) to stay secure.
* **ECDSA (Elliptic Curve):** Uses the "Discrete Logarithm Problem" on an algebraic curve. It provides the same security as RSA but with significantly smaller keys (256 bits vs 3072 bits), making it faster and more efficient for modern cloud environments.

### 2. The Signing Process (Not Encryption!)
A common misconception is that certificates are "encrypted" with a private key. They are actually **Signed**:
1.  **Hashing:** The certificate data (Name, Public Key, Dates) is put through a hash function (SHA-256) to create a unique fingerprint.
2.  **Signing:** The CA uses its **Private Key** to "encrypt" only that fingerprint. 
3.  **Verification:** A client (like a browser) uses the CA's **Public Key** to decrypt the fingerprint. If the browser's calculated fingerprint matches the decrypted one, the certificate is authentic.

### 3. The Chain of Trust
PKI works because of a hierarchy. Your OS or Browser comes pre-installed with "Root" certificates. When you visit a site, the browser checks if that site's certificate was signed by a Root (or an Intermediate trusted by a Root). If the signature chain is unbroken, the green padlock appears.

---

## 🛠 Prerequisites
* Python 3.x
* The `cryptography` library:
    ```bash
    pip install cryptography
    ```

## 🔍 Verification
You can verify your work using the standard `openssl` CLI tool:
```bash
# Check the text content of your certificate
openssl x509 -in server_signed.crt -text -noout

# Verify the chain of trust
openssl verify -CAfile root_cert.pem server_signed.crt
```
