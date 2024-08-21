"use client";

import "./messagebox.css";
import React, { useEffect, useState, useRef } from "react";
import { Message } from "./types";
import Box from "@mui/material/Box";
import RobotIcon from '@mui/icons-material/SmartToy';
import UserIcon from '@mui/icons-material/Person';

interface MessageProps {
    message: Message;
}

const MessageBox: React.FC<MessageProps> = ({message}) => {
    return (
        <>
        <Box whiteSpace="pre-line"
            my={0}
            display="flex"
            gap={2}
            p={2}
            sx={{ width: "95%", border: "2px solid grey", borderRadius: 1 }}
        >
            {message.kind == "sent" ? <UserIcon/> : <RobotIcon/> }
            <p>{message.data}</p>
        </Box>
        <Box>
            <p>moderated: {message.moderated}</p>
        </Box>
        <Box>
            <p>approved: {message.approved}</p>
        </Box>
        </>
    );
};

export default MessageBox;
