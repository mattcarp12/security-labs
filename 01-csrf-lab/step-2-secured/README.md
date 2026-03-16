# Step 2: Secured via CSRF Tokens

This folder demonstrates the **Synchronizer Token Pattern**, the industry-standard cryptographic defense against Cross-Site Request Forgery. 

By requiring a unique, unguessable token on every state-changing request, the server can definitively verify that the request originated from the legitimate frontend, not a hidden attacker script.

## The Defensive Mechanics
1. **Token Generation:** Upon login, the Go backend generates a cryptographically secure 32-byte string and stores it in the database alongside the user's session.
2. **Explicit Fetching:** The frontend JavaScript explicitly requests this token (`GET /csrf-token`).
3. **The Same-Origin Policy (SOP):** If an attacker tries to fetch this token from their malicious domain, the browser's built-in SOP blocks them from reading the response. 
4. **Request Validation:** The frontend attaches the token as a custom header (`X-CSRF-Token`). The backend rejects any transfer request that lacks a matching token.

## How to Start the Lab

1. Open your terminal in this directory (`01-csrf-lab/step-2-secured`).
2. Build and start the Docker containers:
   ```bash
   docker-compose up --build -d
   ```

## How to Verify the Defense

1. Test the Legitimate App
    - Go to `http://localhost:3000` and login `(admin/password123)`
    - Click **Secure Transfer $100**
    - *The javascript fetches the CSRF token and attaches it to the request header. The transfer succeeds*

2. Test the Attacker's Exploit
    - Open a new tab and go to `http://localhost:4000/evil.html`
    - *The malicious form will still automatically submit, and the browser will still attach the session cookie. But because there is no CSRF Token, the request fails*

3. Check the logs:
    - Run `docker-compose logs api` in your terminal
    - You will see the server actively rejected the attacker's request with a 403 Forbidden - Missing CSRF Token error because the attacker's script could not steal the token due to the browser's Same-Origin Policy (SOP).

## Clean Up
```bash
docker-compose down -v
```