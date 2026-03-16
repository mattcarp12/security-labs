# Step 2: Client-Side XSS (Reflected & DOM-Based)

This lab demonstrates two forms of Cross-Site Scripting that do not rely on a database to store the payload: Reflected XSS and DOM-Based XSS.

## 1. Reflected XSS
The payload is included in the URL query parameters (e.g., `?q=<script>...`). The backend server reads this parameter and blindly "reflects" it into the HTML response. When the victim clicks the malicious link, the server hands the payload back to their browser, which executes it.

## 2. DOM-Based XSS
The vulnerability exists entirely within the frontend JavaScript. The script reads malicious data from the URL fragment (the `#hash`), decodes it, and writes it directly to the DOM using dangerous methods like `innerHTML`. The backend server is never involved and never sees the payload.

## How to Test

1. Start the environment: `docker-compose up --build -d`
2. **Test Reflected:** Navigate to `http://localhost:3000` and submit `<script>alert(1)</script>` in the search box.
3. **Test DOM-Based:** Open a new tab and navigate exactly to: 
   `http://localhost:3000/#<img src="x" onerror="alert('DOM XSS')">`
4. Clean up: `docker-compose down`