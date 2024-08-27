"use client";

import React, { useRef, useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import ImageBox from "./image_box";
import { ModRequest } from "./types";
import { Box, Button, Divider, Typography } from "@mui/material";
import { useRouter } from "next/navigation";
import html2pdf from "html2pdf.js";
import "./confirm.css";

const Confirm = () => {
    const router = useRouter();
    const messageUrl = process.env.NEXT_PUBLIC_MESSAGE_REST ?? "undefined";
    const completionUrl =
        process.env.NEXT_PUBLIC_COMPLETION_REST ?? "undefined";
    const pdfUrl = process.env.NEXT_PUBLIC_PDF_REST ?? "undefined";

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

    // not used
    const handleSaveAsPDF = (filepath: string) => {
        if (contentRef.current) {
            const element = contentRef.current;
            const options = {
                // margin: [17, 15], // top-bottom, left-right (mm), bad config
                // margin: [20, 10], // top-bottom, left-right (mm), bad config
                margin: [20, 15], // top-bottom, left-right (mm)
                pagebreak: { mode: ["avoid-all"], before: ".page-break" },
                filename: filepath,
                image: { type: "jpeg", quality: 0.98 },
                html2canvas: {
                    useCORS: true,
                    ignoreElements: (element: HTMLElement) => {
                        // Customize this function to ignore elements based on conditions
                        return element.classList.contains("ignore-pdf");
                    },
                },
                jsPDF: {
                    unit: "mm",
                    format: "a4",
                    orientation: "portrait",
                },
            };

            html2pdf().from(element).set(options).save();
            console.log("%O saved");
        }
    };

    const openEmailClient = () => {
        const subject = `Family Day 2024: ${title}`;
        // const body = `*****************************\n* Attach the generated PDF here*\n******************************\n\n\nThank you very much for joining the family day 2024 and making \"${title}\". I hope you enjoyed the story & image generation!\n\nThank you,\n2024 Family Day`;
        const body =
            (process.env.NEXT_PUBLIC_EMAIL_HEADER ?? "") +
            `\n\n\n*****************************\n* Attach the generated PDF here*\n******************************\n\n\n` +
            (process.env.NEXT_PUBLIC_EMAIL_FOOTER ?? "");

        const mailtoLink = `mailto:${email}?subject=${encodeURIComponent(
            subject
        )}&body=${encodeURIComponent(body)}`;

        if (typeof window !== "undefined") {
            window.location.href = mailtoLink;
        }
    };

    const registerInMongoDB = async (): Promise<boolean> => {
        console.log("registerInMongoDB");

        let params: string =
            "?msgid_txt=" +
            encodeURIComponent(msgid_txt ?? "") +
            "&msgid_url=" +
            encodeURIComponent(msgid_url ?? "") +
            "&email=" +
            encodeURIComponent(email ?? "") +
            "&title=" +
            encodeURIComponent(title ?? "") +
            "&user=" +
            encodeURIComponent(user ?? "") +
            "&employee=" +
            encodeURIComponent(employee ?? "");
        const encoded_url = completionUrl + params;
        console.log(encoded_url);

        const response = await fetch(encoded_url);
        if (!response.ok) {
            console.error("failed to fetch %O", encoded_url);
            return false;
        }
        return true;
    };

    const downloadPDF = async (): Promise<boolean> => {
        console.log("downloadPDF");

        let params: string =
            "?msgid_txt=" +
            encodeURIComponent(msgid_txt ?? "") +
            "&msgid_url=" +
            encodeURIComponent(msgid_url ?? "") +
            "&email=" +
            encodeURIComponent(email ?? "") +
            "&title=" +
            encodeURIComponent(title ?? "") +
            "&user=" +
            encodeURIComponent(user ?? "");

        const encoded_url = pdfUrl + params;
        console.log("* downloading %O", encoded_url);

        try {
            const response = await fetch(encoded_url, {
                method: "GET",
                headers: {
                    Accept: "application/octet-stream",
                },
            });
            if (!response.ok) {
                console.error("failed to fetch %O", encoded_url);
                return false;
            }

            // Convert response to a Blob and create a download link
            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const link = document.createElement("a");
            link.href = url;
            link.download = email + ".pdf";
            document.body.appendChild(link);
            link.click();
            link.remove();
            return true;
        } catch (error) {
            console.error("There was an error downloading the file:", error);
            return false;
        }
    };

    const onComplete = async () => {
        // This is more beautiful, but can't automate
        // window.print();

        // PDF
        const invalidChars = /[\/:*?"<>|\\]/g;

        // This generates a PDF only by Javascript, but characters can be
        // torn up and down between pages.
        // const filename = (email + "_" + title).replace(invalidChars, "_");
        // handleSaveAsPDF(`${filename}.pdf`);

        // generate PDF
        const pdfResult: boolean = await downloadPDF();
        if (!pdfResult) {
            alert("Failed to download a PDF. Please try again");
            return;
        }

        // send email
        openEmailClient();

        // register it in mongodb
        const mongoResult: boolean = await registerInMongoDB();
        if (!mongoResult) {
            alert("Failed to register the completion record. Please try again");
        }
    };

    useEffect(() => {
        // get msg
        fetch(messageUrl + "?msgid=" + msgid_txt)
            .then((res) => res.text())
            .then((data) => {
                setText(data);
            })
            .catch((error) => console.error("Error fetching data:", error));
        // get url
        fetch(messageUrl + "?msgid=" + msgid_url)
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
                        border: "none",
                    }}
                >
                    <Box>
                        <Box
                            p={1}
                            sx={{
                                display: "flex",
                                justifyContent: "center",
                                alignItems: "center",
                                textAlign: "center",
                            }}
                        >
                            <Typography
                                variant="h3"
                                color="lightgray"
                                sx={{ fontSize: "1rem" }}
                            >
                                [ Please review, scroll down to the bottom and
                                generate a PDF ]
                            </Typography>
                        </Box>

                        <Box
                            p={2}
                            sx={{
                                display: "flex",
                                justifyContent: "center",
                                alignItems: "center",
                                textAlign: "center",
                            }}
                        >
                            <Typography variant="h2" sx={{ fontSize: "2rem" }}>
                                {title}
                            </Typography>
                        </Box>
                        <Divider />

                        <Box
                            p={1}
                            sx={{
                                display: "flex",
                                justifyContent: "center",
                                alignItems: "center",
                                textAlign: "center",
                            }}
                        >
                            <Typography
                                variant="h2"
                                sx={{ fontSize: "1.5rem" }}
                            >
                                {user}
                            </Typography>
                        </Box>
                        <Box
                            p={1}
                            sx={{
                                gap: 2,
                                display: "flex",
                                justifyContent: "flex-end",
                            }}
                        >
                            Employee: <b>{employee}</b>
                        </Box>
                        <Box
                            p={1}
                            sx={{
                                gap: 2,
                                display: "flex",
                                justifyContent: "flex-end",
                            }}
                        >
                            Employee Personal Email: <b>{email}</b>
                        </Box>
                        <Box
                            p={1}
                            sx={{
                                display: "flex",
                                justifyContent: "center",
                                alignItems: "center",
                                textAlign: "center",
                            }}
                        >
                            <ImageBox msg={msg} cn="" />
                        </Box>
                        <div className="page-break"></div>
                        <Box
                            p={2}
                            sx={{
                                gap: 2,
                                display: "flex",
                            }}
                        >
                            {text}
                        </Box>
                    </Box>
                </Box>
                <Box
                    className="ignore-pdf"
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
                        Generate PDF and Send E-mail
                    </Button>
                </Box>
            </div>
        </>
    );
};

export default Confirm;
