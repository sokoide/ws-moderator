"use client";

import "./messagebox.css";
import React, { useEffect, useState, useRef, ChangeEvent } from "react";
import ImageBox from "./image_box";
import { Message, ModRequest } from "./types";
import Box from "@mui/material/Box";
import Checkbox from "@mui/material/Checkbox";
import Divider from "@mui/material/Divider";
import RobotIcon from "@mui/icons-material/SmartToy";
import UserIcon from "@mui/icons-material/Person";

interface MessageProps {
    msg: ModRequest;
    checked: boolean;
    onCheckboxChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
}

const MessageBox: React.FC<MessageProps> = ({
    msg,
    checked,
    onCheckboxChange,
}) => {
    const getHighlightedText = (text: string) => {
        const parts = text.split(/(\/imagine)/g);
        return parts.map((part, index) =>
            part === "/imagine" ? (
                <span key={index} className="highlight">
                    {part}
                </span>
            ) : (
                part
            )
        );
    };

    return (
        <>
            <div className="zoom">
                <Box
                    whiteSpace="pre-line"
                    my={0}
                    display="flex"
                    gap={2}
                    p={2}
                    sx={{
                        width: "95%",
                        border: "2px solid grey",
                        borderRadius: 1,
                    }}
                >
                    {msg.client_id === "bot" ? (
                        <Checkbox checked={checked} onChange={onCheckboxChange} />
                    ) : (
                        ""
                    )}
                    {msg.client_id !== "bot" ? <UserIcon /> : <RobotIcon />}
                    <Box display="flex" flexDirection="column">
                        {msg.message.kind === "url" ? (
                            <ImageBox msg={msg} />
                        ) : (
                            <div className="box">
                                <p>{getHighlightedText(msg.message.data)}</p>
                            </div>
                        )}
                        <Divider />
                        <Box sx={{ color: "lightgrey" }}>
                            <p>
                                msgid: {msg.id}, kind: {msg.message.kind}
                            </p>
                        </Box>
                    </Box>
                </Box>
            </div>
            <br />
        </>
    );
};

export default MessageBox;
