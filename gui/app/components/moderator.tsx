"use client";

import "./chat.css";
import React, { useEffect, useState, useRef } from "react";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import MessageBox from "./messagebox";
import { Message } from "./types";
import { green } from "@mui/material/colors";

const Moderator = () => {
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const [inputValue, setInputValue] = useState("");
    const [messages, setMessages] = useState<Message[]>([]);
    const messagePaneRef = useRef<HTMLDivElement>(null);
    const taRef = useRef(null);

    useEffect(() => {
        const ws = new WebSocket("ws://127.0.0.1:80/moderator");

        ws.onopen = () => {
            console.log("useEffect: onopen");
        };

        ws.onmessage = (e) => {
            console.log("useEffect: onmessage: %O", e.data);
            setMessages((prevMessages) => [
                ...prevMessages,
                { kind: "received", data: e.data },
            ]);
        };

        ws.onerror = (error) => {
            console.error("WebSocket error: %O", error);
        };

        ws.onclose = () => {
            console.log("WebSocket connection closed");
        };

        setSocket(ws);

        return () => {
            console.log("useEffect: return");
            ws.close();
        };
    }, []);

    useEffect(() => {
        if (messagePaneRef.current) {
            messagePaneRef.current.scrollTop = messagePaneRef.current.scrollHeight;
        }
    }, [messages]);

    const onSend = () => {
        console.log("onSend: %O", inputValue);
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(inputValue);
            setMessages((prevMessages) => [
                ...prevMessages,
                { kind: "sent", data: inputValue },
            ]);
            setInputValue("");
        } else {
            console.error("WebSocket is not open");
        }
    };

    const handleInputChange = (e) => {
        setInputValue(e.target.value);
    };

    let userText = "1st line\n2nd line\n3rd line";


    return (
        <>
            <Box
                sx={{
                    display: "flex",
                    flexDirection: "column",
                    height: "100vh",
                }}
            >
                <Box ref={messagePaneRef} sx={{ height: "85%", padding: 2, overflowY: 'auto', overflowX: 'auto' }}>
                    {messages.map(function (msg, i) {
                        return <MessageBox message={msg} />;
                    })}
                </Box>

                <Box
                    sx={{
                        padding: 2,
                        backgroundColor: green[300],
                        textAlign: "center",
                    }}
                >
                    <textarea
                        ref={taRef}
                        value={inputValue}
                        onChange={handleInputChange}
                        placeholder="Type here..."
                        rows={4}
                        className="bordered-input"
                    />

                    <Button variant="contained" onClick={onSend}>
                        Send
                    </Button>
                </Box>
            </Box>
        </>
    );
};

export default Moderator;
