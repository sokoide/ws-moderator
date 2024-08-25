"use client";

import React, { useRef, useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import ImageBox from "../components/image_box";
import { Message, ModRequest } from "../components/types";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import { useRouter } from "next/navigation";
import Divider from "@mui/material/Divider";
import html2pdf from "html2pdf.js";

// import dynamic from "next/dynamic";

const ConfirmPage = () => {
    const router = useRouter();
    const restUrl = process.env.NEXT_PUBLIC_MESSAGE_REST ?? "undefined";

    const sp = useSearchParams();
    const msgid_txt = sp.get("msgid_txt");
    const msgid_url = sp.get("msgid_url");
    const title = sp.get("title");
    const user = sp.get("user");
    const employee = sp.get("employee");
    const email = sp.get("email");

    const [text, setText] = useState<string>("");
    const [url, setUrl] = useState<string>("");
    const [msg, setMsg] = useState<ModRequest>({
        id: "",
        client_id: "",
        user_email: "",
        message: { kind: "url", data: "TBD" },
        approved: false,
        moderated: false,
    });
    const contentRef = useRef<HTMLDivElement>(null);

    const onBack = () => {
        router.push("/");
    };

    const handleSaveAsPDF = (filepath: string) => {
        if (contentRef.current) {
            const element = contentRef.current;
            const options = {
                margin: [0.5, 0.5],
                filename: filepath,
                html2canvas: { scale: 2, useCORS: true },
                jsPDF: {
                    unit: "in",
                    format: "A4",
                    orientation: "portrait",
                },
            };

            html2pdf().from(element).set(options).save();
            console.log("%O saved");
        }
    };

    const openEmailClient = () => {
        const subject = `Family Day 2024: ${title}`;
        const body = `*****************************\n* Attach the generated PDF here*\n******************************\n\n\nThank you very much for joining the family day 2024 and making \"${title}\". I hope you enjoyed the story & image generation!\n\nThank you,\n2024 Family Day`;
        const mailtoLink = `mailto:${email}?subject=${encodeURIComponent(
            subject
        )}&body=${encodeURIComponent(body)}`;

        if (typeof window !== "undefined") {
            window.location.href = mailtoLink;
        }
    };

    const onComplete = () => {
        // TODO: register it in mongodb

        // PDF
        const invalidChars = /[\/:*?"<>|\\]/g;

        const filename = (email + "_" + title).replace(invalidChars, "_");
        handleSaveAsPDF(`${filename}.pdf`);

        // send email
        openEmailClient();
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
                    id: "",
                    client_id: "",
                    user_email: "",
                    message: {
                        kind: "url",
                        data: data,
                    },
                    approved: false,
                    moderated: false,
                });
            })
            .catch((error) => console.error("Error fetching data:", error));
    }, []);

    return (
        <>
            <div ref={contentRef}>
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
                                Title: <b>{title}</b>
                            </p>
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
                            <ImageBox msg={msg} cn="" />
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
                        Send E-mail and Complete
                    </Button>
                </Box>
            </div>
        </>
    );
};

export default ConfirmPage;
