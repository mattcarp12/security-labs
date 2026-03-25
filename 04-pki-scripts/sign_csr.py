import datetime
from cryptography import x509
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives import serialization

def sign_certificate_request():
    print("--- Lab 4: Signing the CSR to Issue a Certificate ---")

    # 1. Load the Root CA Private Key (The "Stamp")
    try:
        with open("private_key.pem", "rb") as f:
            root_private_key = serialization.load_pem_private_key(f.read(), password=None)
    except FileNotFoundError:
        print("Error: Root private_key.pem not found.")
        return

    # 2. Load the Root CA Certificate (To get the Issuer Name)
    try:
        with open("root_cert.pem", "rb") as f:
            root_cert = x509.load_pem_x509_certificate(f.read())
    except FileNotFoundError:
        print("Error: root_cert.pem not found.")
        return

    # 3. Load the CSR (The "Application")
    try:
        with open("server_request.csr", "rb") as f:
            csr = x509.load_pem_x509_csr(f.read())
    except FileNotFoundError:
        print("Error: server_request.csr not found.")
        return

    # 4. Sign the CSR to create a "Leaf" Certificate
    # We take the data from the CSR and "stamp" it with the Root CA's authority.
    cert = x509.CertificateBuilder().subject_name(
        csr.subject
    ).issuer_name(
        root_cert.subject # The Root CA is the "Issuer"
    ).public_key(
        csr.public_key() # We pull the Server's Public Key from the CSR
    ).serial_number(
        x509.random_serial_number()
    ).not_valid_before(
        datetime.datetime.utcnow()
    ).not_valid_after(
        # Standard leaf certificates are usually valid for 1 year (365 days)
        datetime.datetime.utcnow() + datetime.timedelta(days=365)
    ).sign(root_private_key, hashes.SHA256())

    # 5. Save the final Signed Certificate
    with open("server_signed.crt", "wb") as f:
        f.write(cert.public_bytes(serialization.Encoding.PEM))

    print("\nSuccess! Issued certificate: server_signed.crt")
    print("This certificate is now 'vouched for' by your Root CA.")

if __name__ == "__main__":
    sign_certificate_request()