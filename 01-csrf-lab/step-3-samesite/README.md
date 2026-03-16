# Step 3: Secured via SameSite Cookies

This folder demonstrates how the `SameSite` cookie attribute provides a robust, browser-level defense against Cross-Site Request Forgery. 

To prove the effectiveness of this defense, the CSRF Token validation from Step 2 has been **intentionally disabled** in this step's backend code. This isolates the variable and shows that `SameSite=Strict` can defeat the attack entirely on its own.

## The Defensive Mechanics
When the Go server issues the session cookie upon login, it attaches the `SameSite=Strict` directive. This acts as a strict set of instructions for the user's web browser:

1. The browser checks the site (the domain currently in the user's address bar).
2. The browser checks the destination of the outgoing request.
3. If the origin and destination do not match (a cross-site request), the browser **refuses to attach the cookie**.
NOTE: The **origin** and **site** are two different concepts. The **origin** consists of the scheme, domain, and port. The **site** consists only of the scheme and domain. Thus, in step 2 below, you need to use `127.0.0.1:4000` instead of `localhost:4000` for the attacker's site, otherwise they will appear to be the same site and the browser will attach the cookie!

Because the attacker's hidden form submission originates from a different domain, the browser strips the session cookie before the request leaves the victim's computer. The server receives an unauthenticated request and safely rejects it.


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
    - *It will succeed because the request originated from the same site, so the browser attached the cookie.*

2. Test the Attacker's Exploit
    - Open a new tab and go to `http://127.0.0.1:4000/evil.html`
    - *The hidden form will submit, but the server will reject it with a `401 Unauthorized` error because the browser withheld the session cookie.*

3. Check the logs:
    - Run `docker-compose logs api` in your terminal
    - You will see the server actively rejected the attacker's request with a 403 Forbidden - Missing CSRF Token error because the attacker's script could not steal the token due to the browser's Same-Origin Policy (SOP).

## Clean Up
```bash
docker-compose down -v
```