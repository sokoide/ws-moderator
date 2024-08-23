"use client";

import "./chat.css";
import React, { useContext, useEffect, useState, useRef } from "react";
import { AppContext } from "@/context/app-context";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Divider from "@mui/material/Divider";
import MessageBox from "./messagebox";
import { Message, ModRequest } from "./types";
import { blue } from "@mui/material/colors";
import { v4 as uuid } from "uuid";

const Chat = () => {
    const { loginInfo, login, logout } = useContext(AppContext);
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const [inputValue, setInputValue] = useState("");
    const [messages, setMessages] = useState<ModRequest[]>([]);
    const messagePaneRef = useRef<HTMLDivElement>(null);
    const taRef = useRef(null);

    if (loginInfo.loggedIn === false) {
        window.location.href = "/login";
    }

    useEffect(() => {
        // TODO: replace it with socket.io?
        const ws = new WebSocket("ws://127.0.0.1:80/chat");

        ws.onopen = () => {
            console.log("useEffect: onopen");
        };

        ws.onmessage = (e) => {
            let msg = JSON.parse(e.data) as ModRequest
            console.log("useEffect: onmessage: %O", msg);
            setMessages((prevMessages) => [
                ...prevMessages,
                msg,
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
            let msg: ModRequest = {
                id: "",
                client_id: "",
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

    let userText = "1st line\n2nd line\n3rd line";

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
                        height: "85%",
                        padding: 2,
                        overflowY: "auto",
                        overflowX: "auto",
                    }}
                >
                    {messages.map(function (msg, i) {
                        return <MessageBox msg={msg} moderator={false} key={uuid()} />;
                    })}
                </Box>

                <Box
                    sx={{
                        padding: 2,
                        backgroundColor: blue[50],
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
        </div>
    );
};

export default Chat;
