"use client";

import "./chat.css";
import React, { useContext, useEffect, useState, useRef } from "react";
import { AppContext } from "@/context/app-context";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import MessageBox from "./messagebox";
import { ChatItem } from "./types";
import { green } from "@mui/material/colors";

const Moderator = () => {
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const [inputValue, setInputValue] = useState("");
    const [messages, setMessages] = useState<ChatItem[]>([]);
    const messagePaneRef = useRef<HTMLDivElement>(null);
    const taRef = useRef(null);
    const { loginInfo, login, logout } = useContext(AppContext);

    useEffect(() => {
        const ws = new WebSocket("ws://127.0.0.1:80/moderator");

        ws.onopen = () => {
            console.log("useEffect: onopen");
        };

        ws.onmessage = (e) => {
            console.log("useEffect: onmessage: %O", e.data);
            setMessages((prevMessages) => [
                ...prevMessages,
                { direction: "received", data: e.data },
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
            messagePaneRef.current.scrollTop =
                messagePaneRef.current.scrollHeight;
        }
    }, [messages]);

    const onSend = () => {
        console.log("onSend: %O", inputValue);
        if (socket && socket.readyState === WebSocket.OPEN) {
            let message: ChatItem = {
                direction: "sent",
                data: inputValue,
                moderated: false,
                approved: false,
            };
            socket.send(JSON.stringify(message));
            setMessages((prevMessages) => [
                ...prevMessages,
                message,
            ]);
            setInputValue("");
        } else {
            console.error("WebSocket is not open");
        }
    };

    const handleInputChange = (e) => {
        setInputValue(e.target.value);
    };

    return (
        <div>
            <Box
                sx={{
                    display: "flex",
                    flexDirection: "column",
                    height: "95vh",
                }}
            >
                <Box
                    ref={messagePaneRef}
                    sx={{
                        height: "100%",
                        padding: 2,
                        overflowY: "auto",
                        overflowX: "auto",
                    }}
                >
                    {messages.map(function (msg, i) {
                        return <MessageBox message={msg} />;
                    })}
                </Box>
            </Box>
        </div>
    );
};

export default Moderator;
