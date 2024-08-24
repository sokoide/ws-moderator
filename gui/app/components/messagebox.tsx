"use client";

import "./messagebox.css";
import React, { useEffect, useState, useRef } from "react";
import { Message, ModRequest } from "./types";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import RobotIcon from "@mui/icons-material/SmartToy";
import UserIcon from "@mui/icons-material/Person";

interface MessageProps {
    msg: ModRequest;
}

const MessageBox: React.FC<MessageProps> = ({ msg }) => {
    // console.info("message: %O", msg);

    return (
        <>
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
                {msg.client_id !== "bot" ? <UserIcon /> : <RobotIcon />}
                <Box flexDirection="column">
                    <p>{msg.message.data}</p>
                    <Divider />
                    <Box sx={{color: "lightgrey"}}>
                        <p>msgid: {msg.id}</p>
                    </Box>
                </Box>
            </Box>
            <br />
        </>
    );
};

export default MessageBox;
