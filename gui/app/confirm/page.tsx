"use client";

import React from "react";

import dynamic from "next/dynamic";
const Confirm = dynamic(() => import("../components/confirm"),{ ssr: false, });


const ConfirmPage = () => {

    return (
        <>
        <Confirm/>
        </>
    );
};

export default ConfirmPage;
