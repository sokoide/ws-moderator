"use client";
import Chat from "./components/chat";
import { useRouter } from "next/navigation";

const HomePage = () => {
    const router = useRouter();

    return (
        <>
            <Chat />
        </>
    );
};

export default HomePage;
