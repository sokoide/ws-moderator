"use client";

import "./messagebox.css";
import React, { MouseEvent } from "react";
import ImageBox from "./image_box";
import { ModRequest } from "./types";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Divider from "@mui/material/Divider";
import RobotIcon from "@mui/icons-material/SmartToy";
import UserIcon from "@mui/icons-material/Person";

interface ModerateButtonsProps {
    msg: ModRequest;
    onApprove: (value: string) => void;
    onDeny: (value: string) => void;
}

const ModerateButtons: React.FC<ModerateButtonsProps> = ({
    msg,
    onApprove,
    onDeny,
}) => {
    const onChildApprove = (e: MouseEvent<HTMLButtonElement>) => {
        console.log("onChildApproveButton: %O", e.currentTarget.value);
        onApprove(e.currentTarget.value);
    };

    const onChildDeny = (e: MouseEvent<HTMLButtonElement>) => {
        console.log("onChildDenyButton: %O", e.currentTarget.value);
        onDeny(e.currentTarget.value);
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
                    <Button
                        variant="contained"
                        value={msg.id}
                        onClick={onChildApprove}
                    >
                        Approve
                    </Button>
                    &nbsp;
                    <Button
                        variant="contained"
                        value={msg.id}
                        onClick={onChildDeny}
                    >
                        Deny
                    </Button>
                </Box>
            </Box>
        </>
    );
};

const ModeratorMessageBox: React.FC<ModerateButtonsProps> = ({
    msg,
    onApprove,
    onDeny,
}) => {
    // console.info("message: %O", msg);

    const onChildApprove = (msgid: string) => {
        console.log("onChildApprove: %O", msgid);
        onApprove(msgid);
    };

    const onChildDeny = (msgid: string) => {
        console.log("onChildDeny: %O", msgid);
        onDeny(msgid);
    };

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
                            {msg.message.kind === "url" ? (
                                <ImageBox msg={msg} />
                            ) : (
                                <div className="box">
                                    <p>
                                        {getHighlightedText(msg.message.data)}
                                    </p>
                                </div>
                            )}
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
