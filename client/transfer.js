// client/transfer.js

async function secureTransfer(event) {
    // Stop the default HTML form submission
    event.preventDefault(); 

    const amount = document.getElementById("transfer-amount").value;
    const resultBox = document.getElementById("result-text");

    try {
        // Step 1: Ask the server for the secret handshake (the CSRF token)
        const tokenResp = await fetch("http://localhost:8080/csrf-token", {
            method: "GET",
            credentials: "include" // <-- Important to include cookies for session identification
        });
        
        // If the user isn't logged in, this will fail
        if (!tokenResp.ok) {
            resultBox.innerText = "Error: Please log in first.";
            return;
        }

        const data = await tokenResp.json();
        const csrfToken = data.csrf;

        // Step 2: Send the transfer request WITH the token in the header
        const transferResp = await fetch(`http://localhost:8080/transfer?amount=${amount}`, {
            method: "POST",
            credentials: "include", // <-- Important to include cookies for session identification
            headers: {
                "X-CSRF-Token": csrfToken // <-- The defensive measure
            }
        });

        if (transferResp.ok) {
            const result = await transferResp.text();
            resultBox.innerText = result;
        } else {
            resultBox.innerText = `Transfer Failed: ${transferResp.statusText}`;
        }

    } catch (error) {
        resultBox.innerText = "Network error occurred.";
    }
}