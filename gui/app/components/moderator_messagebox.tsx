"use client";

import "./messagebox.css";
import React, { useEffect, useState, useRef } from "react";
import { Message, ModRequest } from "./types";
import Moderator from "./moderator";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Divider from "@mui/material/Divider";
import RobotIcon from "@mui/icons-material/SmartToy";
import UserIcon from "@mui/icons-material/Person";

// interface MessageProps {
//     msg: ModRequest;
// }

interface ModerateButtonsProps {
    msg: ModRequest;
    onApprove: (value: string) => void;
    onDeny: (value: string) => void;
}

const ModerateButtons: React.FC<ModerateButtonsProps> = ({ msg, onApprove, onDeny }) => {
    const onChildApprove = (e) => {
        console.log("onChildApproveButton: %O", e.target.value);
        onApprove(e.target.value);
    };

    const onChildDeny = (e) => {
        console.log("onChildDenyButton: %O", e.target.value);
        onDeny(e.target.value);
    };

    return (
        <>
            <Box flexDirection={"column"}>
                <Divider />
                <Box>
                    <p>email: {msg.user_email}</p>
                    <p>msgid: {msg.id}</p>
                </Box>
                <Box flexDirection={"row"}>
                    <Button variant="contained" value={msg.id} onClick={onChildApprove}>
                        Approve
                    </Button>
                    &nbsp;
                    <Button variant="contained" value={msg.id} onClick={onChildDeny}>
                        Deny
                    </Button>
                </Box>
            </Box>
        </>
    );
};

const ModeratorMessageBox: React.FC<ModerateButtonsProps> = ({ msg, onApprove, onDeny }) => {
    console.info("message: %O", msg);

    const onChildApprove = (msgid) => {
        console.log("onChildApprove: %O", msgid);
        onApprove(msgid);
    };

    const onChildDeny = (msgid) => {
        console.log("onChildDeny: %O", msgid);
        onDeny(msgid);
    };

    if (msg.moderated) {
        return <></>;
    } else {
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
                    <Box flexDirection={"column"}>
                        <Box>
                            <p>{msg.message.data}</p>
                        </Box>
                        {msg.moderated || msg.approved ? (
                            ""
                        ) : (
                            <ModerateButtons
                                msg={msg}
                                onApprove={onChildApprove}
                                onDeny={onChildDeny}
                            />
                        )}
                    </Box>
                </Box>
                <br />
            </>
        );
    }
};

export default ModeratorMessageBox;
