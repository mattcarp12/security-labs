# Step 1: Vulnerable CSRF Application

This folder contains a web application intentionally left vulnerable to Cross-Site Request Forgery (CSRF). It demonstrates how browsers automatically attach session cookies to cross-origin requests.

## How to Start the Lab

1. Open your terminal in this directory (`01-csrf-lab/step-1-vulnerable`).
2. Build and start the Docker containers:
   ```bash
   docker-compose up --build -d
   ```
3. Wait a few seconds for the database to initialize and the Go server to connect.

## How to Execute the Attack

1. Log into the legitimate site
    
    - Open your browser to `http://localhost:3000`
    - Login with default credentials: `admin/password123`
    - You will see "Login successful. Cookie set!" The browser now has your session cookie. Verify this in your browser's Dev Console

2. Trigger the Exploit
    - Open a new tab in the same browser window
    - Go to `http://localhost:4000/evil.html`
    - Watch the page. After 1 second, a hidden form will automatically submit a $10,000 transfer request to the legitimate API. Because you are logged in, the browser attaches your cookie, and the server processes the forged transfer.

## How to Stop and Clean Up

```bash
docker-compose down -v
```

