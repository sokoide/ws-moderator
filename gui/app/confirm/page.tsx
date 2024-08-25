"use client";

import React, { useCallback, useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import ImageBox from "../components/image_box";
import { Message, ModRequest } from "../components/types";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import { useRouter } from "next/navigation";
import Divider from "@mui/material/Divider";

const ConfirmPage = () => {
    const router = useRouter();
    const restUrl = process.env.NEXT_PUBLIC_MESSAGE_REST ?? "undefined";
    const sp = useSearchParams();
    const msgid_txt = sp.get("msgid_txt");
    const msgid_url = sp.get("msgid_url");
    const user = sp.get("user");
    const employee = sp.get("employee");
    const email = sp.get("email");

    const [text, setText] = useState<string>("");
    const [url, setUrl] = useState<string>("");
    const [msg, setMsg] = useState<ModRequest>({
        message: { kind: "url", data: "TBD" },
    });

    const onBack = () => {
        router.push("/");
    };

    const onComplete = () => {
        // TODO: send email
        alert(
            "TODO: Make a PDF and send an email here. All done. Thank you for joining us!"
        );
        router.push("/login");
    };

    useEffect(() => {
        // get msg
        fetch(restUrl + "?msgid=" + msgid_txt)
            .then((res) => res.text())
            .then((data) => {
                setText(data);
            })
            .catch((error) => console.error("Error fetching data:", error));
        // get url
        fetch(restUrl + "?msgid=" + msgid_url)
            .then((res) => res.text())
            .then((data) => {
                setUrl(data);
                setMsg({
                    message: {
                        kind: "url",
                        data: data,
                    },
                });
            })
            .catch((error) => console.error("Error fetching data:", error));
    }, []);

    return (
        <>
            <Box
                whiteSpace="pre-line"
                my={4}
                mx={4}
                width="100%"
                display="flex"
                flexDirection="column"
                gap={2}
                p={2}
                sx={{
                    width: "95%",
                    border: "2px solid grey",
                    borderRadius: 1,
                }}
            >
                <Box>
                    <Box p={2}>
                        <p>
                            User: <b>{user}</b>
                        </p>
                        <p>
                            Employee: <b>{employee}</b>
                        </p>
                        <p>
                            Employee Email: <b>{email}</b>
                        </p>
                    </Box>
                    <Box p={2}>
                        <ImageBox msg={msg} />
                    </Box>
                    <Divider />
                    <Box p={2}>
                        <p>{text}</p>
                    </Box>
                </Box>
            </Box>
            <Box
                display="flex"
                gap={2}
                justifyContent="center"
                padding={2}
                flexDirection="row"
            >
                <Button variant="contained" onClick={onBack}>
                    Go Back
                </Button>
                &nbsp;
                <Button variant="contained" onClick={onComplete}>
                    Complete
                </Button>
            </Box>
        </>
    );
};

export default ConfirmPage;
