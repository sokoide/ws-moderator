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
import { TextField, Box, Button, Divider } from "@mui/material";

import MessageBox from "./messagebox";
import { ModRequest } from "./types";
import ClientUtil from "./client_util";
import { v4 as uuid } from "uuid";
import { useRouter } from "next/navigation";

const Chat = () => {
    const router = useRouter();
    const context = useContext(AppContext);
    if (context === null) {
        console.error("context not available");
        return;
    }

    const { loginInfo } = context;
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const [inputValue, setInputValue] = useState("");
    const [messages, setMessages] = useState<ModRequest[]>([]);
    const messagesEndRef = useRef<HTMLDivElement | null>(null);
    const eopRef = useRef<HTMLDivElement | null>(null);
    const titleRef = useRef<HTMLInputElement | null>(null);
    const userRef = useRef<HTMLInputElement | null>(null);
    const employeeRef = useRef<HTMLInputElement | null>(null);

    if (loginInfo != null && loginInfo.loggedIn === false) {
        if (typeof window !== "undefined") window.location.href = "/login";
    }
    const scrollToBottom = () => {
         if (messagesEndRef.current) {
            console.log("*** scrollToBottom");
             messagesEndRef.current.scrollIntoView({ behavior: "smooth" });
         }
    };

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
        setCheckboxStates(
            messages.map((message) => ({
                msg: message,
                checked: false,
            }))
        );

        // Use setTimeout to ensure scroll happens after DOM update
        const timer = setTimeout(() => {
            scrollToBottom();
        }, 10);

        return () => clearTimeout(timer);
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
        let email: string = "";

        if (
            !userRef.current ||
            !employeeRef.current ||
            !titleRef.current ||
            userRef.current.value.trim() === "" ||
            employeeRef.current.value.trim() === "" ||
            titleRef.current.value.trim() === ""
        ) {
            alert("please input Title, User and Employee");
            return;
        }

        checkboxStates.map(({ msg, checked }) => {
            console.log("* %O:%O -> %O", msg.id, msg.message.kind, checked);
            if (checked) {
                if (msg.message.kind === "txt") {
                    id_txt_count++;
                    id_txt = msg.id;
                    email = msg.user_email;
                } else if (msg.message.kind === "url") {
                    id_url_count++;
                    id_url = msg.id;
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

        let target_page: string = "/confirm?";
        let params: string =
            "msgid_txt=" +
            encodeURIComponent(id_txt) +
            "&msgid_url=" +
            encodeURIComponent(id_url) +
            "&email=" +
            encodeURIComponent(email) +
            "&title=" +
            encodeURIComponent(titleRef.current.value) +
            "&user=" +
            encodeURIComponent(userRef.current.value) +
            "&employee=" +
            encodeURIComponent(employeeRef.current.value);
        const encoded_page = target_page + params;
        console.log(encoded_page);
        router.push(encoded_page);
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
        <>
            <div ref={eopRef} className="container">
                <Box>
                    <Box
                        my={4}
                        mx={8}
                        width="100%"
                        sx={{
                            display: "flex",
                            flexDirection: "column",
                            height: "90vh",
                            overflowY: "auto",
                        }}
                    >
                        <Box>
                            {checkboxStates.map(({ msg, checked }) => (
                                <MessageBox
                                    msg={msg}
                                    checked={checked}
                                    onCheckboxChange={handleCheckboxChange(
                                        msg.id
                                    )}
                                    key={uuid()}
                                />
                            ))}
                        </Box>
                        <div ref={messagesEndRef} />
                    </Box>

                    <p>
                        [*] When completed, input <b>Title, User, Employee</b>,
                        check <b>1 text &amp; 1 image</b> →{" "}
                        <b>[Review &amp; Complete]</b> button{" "}
                    </p>
                    <Box
                        display="flex"
                        flexDirection="row"
                        gap={1}
                    >
                        <textarea
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
                    <Box display="flex"
                        justifyContent="center"
                        p={1}
                        gap={4}
                    >
                        <TextField
                            id="title"
                            inputRef={titleRef}
                            label="Title of a book"
                            variant="outlined"
                            style={{ width: "40%" }}
                        />
                        <TextField
                            id="user"
                            inputRef={userRef}
                            label="User who used AI today"
                            variant="outlined"
                        />
                        <TextField
                            id="employee"
                            inputRef={employeeRef}
                            label="Employee Name"
                            variant="outlined"
                        />
                    </Box>
                    <Box
                        display="flex"
                        justifyContent="center"
                        p={1}
                        gap={4}
                    >
                        <Button variant="outlined" onClick={onUncheckAll}>
                            Uncheck All
                        </Button>
                        <Button variant="outlined" onClick={onComplete}>
                            Generate PDF
                        </Button>
                    </Box>
                </Box>
            </div>
        </>
    );
};

export default Chat;
