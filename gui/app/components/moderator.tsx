"use client";

import "./chat.css";
import React, { useContext, useEffect, useState, useRef } from "react";
import { AppContext } from "@/context/app-context";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import ModeratorMessageBox from "./moderator_messagebox";
import { Message, ModRequest } from "./types";
import { green } from "@mui/material/colors";
import { v4 as uuid } from "uuid";

const Moderator = () => {
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const [inputValue, setInputValue] = useState("");
    const [messages, setMessages] = useState<ModRequest[]>([]);
    const messagePaneRef = useRef<HTMLDivElement>(null);
    const taRef = useRef(null);
    const { loginInfo, login, logout } = useContext(AppContext);

    useEffect(() => {
        const ws = new WebSocket("ws://127.0.0.1:80/moderator");

        ws.onopen = () => {
            console.log("useEffect: onopen");
            let msg: ModRequest = {
                id: "",
                client_id: "",
                user_email: "",
                message: {
                    kind: "system",
                    data: "",
                },
                approved: false,
                moderated: false,
            };
            ws.send(JSON.stringify(msg));
        };

        ws.onmessage = (e) => {
            let msg = JSON.parse(e.data) as ModRequest;
            console.log("useEffect: onmessage: %O", msg);
            if (msg.message.kind === "ping") return;
            setMessages((prevMessages) => [...prevMessages, msg]);
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
            let msg: ModRequest = {
                id: "",
                client_id: "",
                user: "user",
                user_email: loginInfo.email,
                message: {
                    kind: "txt",
                    data: inputValue,
                },
                approved: false,
                moderated: false,
            };
            socket.send(JSON.stringify(msg));
            setMessages((prevMessages) => [...prevMessages, msg]);
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
                        return <ModeratorMessageBox msg={msg} key={uuid()} />;
                    })}
                </Box>
            </Box>
        </div>
    );
};

export default Moderator;
