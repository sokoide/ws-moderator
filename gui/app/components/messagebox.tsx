"use client";

import "./messagebox.css";
import React, { useEffect, useState, useRef, ChangeEvent } from "react";
import ImageBox from "./image_box";
import { ModRequest } from "./types";
import { Box, Checkbox, Divider, IconButton } from "@mui/material";
import RobotIcon from "@mui/icons-material/SmartToy";
import UserIcon from "@mui/icons-material/Person";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import { Content } from "next/font/google";

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

    const handleCopy = () => {
        navigator.clipboard
            .writeText(msg.message.data)
            .then(() => {
                alert("Content copied to clipboard!");
            })
            .catch((err) => {
                console.error("Failed to copy: ", err);
            });
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
                        <Checkbox
                            checked={checked}
                            onChange={onCheckboxChange}
                        />
                    ) : (
                        ""
                    )}
                    {msg.client_id !== "bot" ? <UserIcon /> : <RobotIcon />}
                    <Box display="flex" flexDirection="column" width={"100%"}>
                        {msg.message.kind === "url" ? (
                            <ImageBox msg={msg} cn="zoom" />
                        ) : (
                            <div className="box">
                                <p>{getHighlightedText(msg.message.data)}</p>
                            </div>
                        )}
                        <Divider />
                        <Box sx={{color: "lightgrey"}} display="flex" flexDirection="row" justifyContent="space-between">
                            <p>
                                msgid: {msg.id}, kind: {msg.message.kind}
                            </p>
                            <IconButton
                                color="primary"
                                onClick={handleCopy}
                                sx={{ marginLeft: 2 }}
                            >
                                <ContentCopyIcon />
                            </IconButton>
                        </Box>
                    </Box>
                </Box>
            </div>
            <br />
        </>
    );
};

export default MessageBox;
