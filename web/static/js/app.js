const chatArea = document.getElementById("chatArea");
const chatForm = document.getElementById("chatForm");
const messageInput = document.getElementById("messageInput");
const sendBtn = document.getElementById("sendBtn");

let conversationHistory = [];
const MAX_HISTORY = 20;

chatForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const text = messageInput.value.trim();
    if (!text) return;

    appendMessage("user", text);
    conversationHistory.push({ role: "user", content: text });
    if (conversationHistory.length > MAX_HISTORY) {
        conversationHistory = conversationHistory.slice(-MAX_HISTORY);
    }

    messageInput.value = "";
    sendBtn.disabled = true;

    const yodaBubble = appendMessage("yoda", "");
    yodaBubble.classList.add("thinking");
    yodaBubble.textContent = "Yoda is thinking";

    try {
        const response = await fetch("/api/chat/stream", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                message: text,
                history: conversationHistory.slice(0, -1),
            }),
        });

        if (!response.ok) {
            throw new Error(`Server error: ${response.status}`);
        }

        yodaBubble.classList.remove("thinking");
        yodaBubble.textContent = "";

        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        let fullResponse = "";
        let buffer = "";

        while (true) {
            const { done, value } = await reader.read();
            if (done) break;

            buffer += decoder.decode(value, { stream: true });
            const lines = buffer.split("\n");
            buffer = lines.pop();

            for (const line of lines) {
                if (!line.startsWith("data: ")) continue;
                const data = line.slice(6).trim();
                if (data === "[DONE]") continue;

                try {
                    const parsed = JSON.parse(data);
                    if (parsed.error) {
                        yodaBubble.textContent += `\n[Error: ${parsed.error}]`;
                        continue;
                    }
                    if (parsed.token) {
                        fullResponse += parsed.token;
                        yodaBubble.textContent = fullResponse;
                        scrollToBottom();
                    }
                } catch {}
            }
        }

        if (fullResponse) {
            conversationHistory.push({ role: "assistant", content: fullResponse });
        }
    } catch (err) {
        yodaBubble.classList.remove("thinking");
        yodaBubble.textContent = "Disturbed, the Force is. Failed to reach Yoda, the connection has.";
        console.error("Stream error:", err);
    } finally {
        sendBtn.disabled = false;
        messageInput.focus();
    }
});

function appendMessage(role, content) {
    const msg = document.createElement("div");
    msg.className = `message ${role}`;

    if (role === "yoda") {
        const avatar = document.createElement("img");
        avatar.src = "/static/img/yoda-chibi.jpg";
        avatar.alt = "Yoda";
        avatar.className = "avatar";
        msg.appendChild(avatar);
    }

    const bubble = document.createElement("div");
    bubble.className = "bubble";
    bubble.textContent = content;
    msg.appendChild(bubble);

    chatArea.appendChild(msg);
    scrollToBottom();
    return bubble;
}

function scrollToBottom() {
    chatArea.scrollTop = chatArea.scrollHeight;
}
