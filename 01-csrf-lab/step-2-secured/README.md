# Step 2: Secured with CSRF Tokens

This phase of the lab patches the vulnerability by implementing the **Synchronizer Token Pattern**. 

While the browser still automatically attaches the session cookie, the server now requires a secondary, explicitly provided "secret handshake" (the CSRF token) to process any state-changing requests.

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