"use client";

import "./chat.css";
import React, {
    useContext,
    useEffect,
    useState,
    useRef,
    ChangeEvent,
} from "react";
import { AppContext } from "@/context/app-context";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import MessageBox from "./messagebox";
import { ModRequest } from "./types";
import ClientUtil from "./client_util";
import { v4 as uuid } from "uuid";

const Chat = () => {
    const context = useContext(AppContext);
    if (context === null) {
        console.error("context not available")
        return;
    }

    const { loginInfo } = context;
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const [inputValue, setInputValue] = useState("");
    const [messages, setMessages] = useState<ModRequest[]>([]);
    const messagePaneRef = useRef<HTMLDivElement>(null);
    const taRef = useRef(null);

    if (loginInfo != null && loginInfo.loggedIn === false) {
        window.location.href = "/login";
    }

    useEffect(() => {
        const wsUrl = process.env.NEXT_PUBLIC_CHAT_WS ?? "undefined";
        console.log("wsUrl: %O", wsUrl);
        const ws = new WebSocket(wsUrl);

        ws.onopen = () => {
            console.log("useEffect: onopen");
            ClientUtil.sendMessage(
                ws,
                "",
                "",
                loginInfo?.email ?? "",
                "system",
                "",
                false,
                false
            );
        };

        ws.onmessage = (e) => {
            let msg = JSON.parse(e.data) as ModRequest;
            console.log("useEffect: onmessage: %O", msg);
            if (
                (msg.client_id === "bot" && msg.id === "") ||
                msg.user_email === (loginInfo?.email ?? "")
            ) {
                setMessages((prevMessages) => [...prevMessages, msg]);
            }
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
        let msg = ClientUtil.sendMessage(
            socket,
            "",
            "",
            loginInfo?.email ?? "",
            "txt",
            inputValue,
            false,
            false
        );
        if (msg != null) {
            // Note: don't add the message here.
            // the requested message will be sent back from the server
            // setMessages((prevMessages) => [...prevMessages, msg]);
            setInputValue("");
        } else {
            console.error("WebSocket is not open");
        }
    };

    const handleInputChange = (e: ChangeEvent<HTMLTextAreaElement>) => {
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
                        return <MessageBox msg={msg} key={uuid()} />;
                    })}
                </Box>

                <Box className="message-input">
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
