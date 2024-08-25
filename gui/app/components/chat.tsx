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
        console.error("context not available");
        return;
    }

    const { loginInfo } = context;
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const [inputValue, setInputValue] = useState("");
    const [messages, setMessages] = useState<ModRequest[]>([]);
    const messagePaneRef = useRef<HTMLDivElement>(null);
    const taRef = useRef(null);

    if (loginInfo != null && loginInfo.loggedIn === false) {
        if (typeof window !== "undefined") window.location.href = "/login";
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

        setCheckboxStates(
            messages.map((message) => ({
                msg: message,
                checked: false,
            }))
        );
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

    const onComplete = () => {
        console.log("onComplete");
        let id_txt_count: number = 0;
        let id_url_count: number = 0;
        let id_txt: string = "";
        let id_url: string = "";

        checkboxStates.map(({ msg, checked }) => {
            console.log("* %O:%O -> %O", msg.id, msg.message.kind, checked);
            if (checked) {
                if (msg.message.kind === "txt") {
                    id_txt_count++;
                    id_txt = msg.id;
                } else if (msg.message.kind === "url") {
                    id_url_count++;
                    id_url = msg.user_email;
                }
            }
        });
        if (id_txt_count !== 1) {
            alert("please select only one text");
            return;
        }
        if (id_url_count !== 1) {
            alert("please select only one image");
            return;
        }

        // TODO: complete
    };

    const onUncheckAll = () => {
        console.log("onUncheckAll");
        setCheckboxesState(false);
    };

    const handleInputChange = (e: ChangeEvent<HTMLTextAreaElement>) => {
        setInputValue(e.target.value);
    };

    // State to manage multiple checkboxes
    const [checkboxStates, setCheckboxStates] = useState<
        { msg: ModRequest; checked: boolean }[]
    >([]);

    // Handler to update the checkbox state
    const handleCheckboxChange =
        (id: string) => (event: React.ChangeEvent<HTMLInputElement>) => {
            setCheckboxStates((prevStates) =>
                prevStates.map((checkbox) =>
                    checkbox.msg.id === id
                        ? { ...checkbox, checked: event.target.checked }
                        : checkbox
                )
            );
        };

    // Function to programmatically check/uncheck checkboxes
    const setCheckboxesState = (checked: boolean) => {
        setCheckboxStates((prevStates) =>
            prevStates.map((checkbox) => ({ ...checkbox, checked }))
        );
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
                        height: "85%",
                        padding: 2,
                        overflowY: "auto",
                        overflowX: "auto",
                    }}
                >
                    {checkboxStates.map(({ msg, checked }) => (
                        <MessageBox
                            msg={msg}
                            checked={checked}
                            onCheckboxChange={handleCheckboxChange(msg.id)}
                            key={uuid()}
                        />
                    ))}
                </Box>

                <Box className="message-input" display="flex" flexDirection="row" gap={1}>
                    <textarea
                        ref={taRef}
                        value={inputValue}
                        onChange={handleInputChange}
                        placeholder="Type here..."
                        rows={4}
                        className="bordered-input"
                    />
                    <Button variant="contained" onClick={onSend}>
                        Ask
                    </Button>
                </Box>
                <Box display="flex" gap={2} justifyContent="center" padding={2}>
                    <Button variant="contained" onClick={onUncheckAll}>
                        Uncheck All
                    </Button>
                    &nbsp;
                    <Button variant="contained" onClick={onComplete}>
                        Complete
                    </Button>
                </Box>
            </Box>
        </div>
    );
};

export default Chat;
